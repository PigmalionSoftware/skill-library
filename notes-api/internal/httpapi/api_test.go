package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"notes-api/internal/notes"
)

func testAPI(t *testing.T) (http.Handler, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "notes.json")
	store, err := notes.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	return New(store), path
}

func request(handler http.Handler, method, path, body, contentType string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, path, strings.NewReader(body))
	if contentType != "" {
		r.Header.Set("Content-Type", contentType)
	}
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	return w
}

func TestCRUD(t *testing.T) {
	handler, _ := testAPI(t)
	w := request(handler, "GET", "/notes", "", "")
	if w.Code != http.StatusOK || strings.TrimSpace(w.Body.String()) != "[]" {
		t.Fatalf("empty list: %d %s", w.Code, w.Body.String())
	}
	w = request(handler, "POST", "/notes", `{"title":"  Shopping  ","content":"Milk"}`, "application/json; charset=utf-8")
	if w.Code != http.StatusCreated || w.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("create: %d %s", w.Code, w.Body.String())
	}
	var created notes.Note
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	path := "/notes/" + created.ID
	if w.Header().Get("Location") != path || created.Title != "Shopping" || created.Content != "Milk" {
		t.Fatalf("unexpected create response: %+v", created)
	}
	w = request(handler, "GET", path, "", "")
	var fetched notes.Note
	if err := json.Unmarshal(w.Body.Bytes(), &fetched); err != nil || w.Code != http.StatusOK || fetched != created {
		t.Fatalf("get: %d %s, %v", w.Code, w.Body.String(), err)
	}
	w = request(handler, "PUT", path, `{"title":"Updated","content":"Bread"}`, "application/json")
	var updated notes.Note
	if err := json.Unmarshal(w.Body.Bytes(), &updated); err != nil || w.Code != http.StatusOK {
		t.Fatalf("update: %d %s, %v", w.Code, w.Body.String(), err)
	}
	if updated.ID != created.ID || updated.CreatedAt != created.CreatedAt || updated.Title != "Updated" || updated.Content != "Bread" {
		t.Fatalf("invalid updated note: %+v", updated)
	}
	w = request(handler, "GET", "/notes", "", "")
	var listed []notes.Note
	if err := json.Unmarshal(w.Body.Bytes(), &listed); err != nil || w.Code != http.StatusOK || len(listed) != 1 || listed[0] != updated {
		t.Fatalf("list: %d %s, %v", w.Code, w.Body.String(), err)
	}
	// PUT replaces all editable fields, so omitted content becomes empty.
	w = request(handler, "PUT", path, `{"title":"Title only"}`, "application/json")
	if err := json.Unmarshal(w.Body.Bytes(), &updated); err != nil || w.Code != http.StatusOK || updated.Content != "" {
		t.Fatalf("replace: %d %s, %v", w.Code, w.Body.String(), err)
	}
	w = request(handler, "DELETE", path, "", "")
	if w.Code != http.StatusNoContent || w.Body.Len() != 0 {
		t.Fatalf("delete: %d %s", w.Code, w.Body.String())
	}
	for _, method := range []string{"GET", "PUT", "DELETE"} {
		w = request(handler, method, path, `{"title":"Missing"}`, "application/json")
		if w.Code != http.StatusNotFound {
			t.Errorf("missing %s: %d %s", method, w.Code, w.Body.String())
		}
	}
}

func TestInvalidRequests(t *testing.T) {
	handler, _ := testAPI(t)
	cases := []struct {
		name, body, contentType string
		status                  int
	}{
		{"empty", "", "application/json", 400},
		{"malformed", "{", "application/json", 400},
		{"null", "null", "application/json", 400},
		{"array", "[]", "application/json", 400},
		{"blank title", `{"title":"   "}`, "application/json", 400},
		{"unknown field", `{"title":"Note","extra":true}`, "application/json", 400},
		{"wrong type", `{"title":123}`, "application/json", 400},
		{"multiple objects", `{"title":"One"} {"title":"Two"}`, "application/json", 400},
		{"trailing garbage", `{"title":"One"} garbage`, "application/json", 400},
		{"long title", `{"title":"` + strings.Repeat("a", 201) + `"}`, "application/json", 400},
		{"long content", `{"title":"Note","content":"` + strings.Repeat("a", 10001) + `"}`, "application/json", 400},
		{"oversize", `{"title":"Note","content":"` + strings.Repeat("a", maxBodyBytes) + `"}`, "application/json", 413},
		{"oversize trailing whitespace", `{"title":"Note"}` + strings.Repeat(" ", maxBodyBytes), "application/json", 413},
		{"missing media type", `{"title":"Note"}`, "", 415},
		{"wrong media type", `{"title":"Note"}`, "text/plain", 415},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := request(handler, "POST", "/notes", tc.body, tc.contentType)
			if w.Code != tc.status || w.Header().Get("Content-Type") != "application/json" {
				t.Fatalf("got %d %s, want %d", w.Code, w.Body.String(), tc.status)
			}
			var failure map[string]string
			if err := json.Unmarshal(w.Body.Bytes(), &failure); err != nil || failure["error"] == "" {
				t.Fatalf("invalid error response: %s", w.Body.String())
			}
		})
	}
	w := request(handler, "GET", "/notes", "", "")
	if strings.TrimSpace(w.Body.String()) != "[]" {
		t.Fatal("invalid requests created notes")
	}
}

func TestRoutingAndHealth(t *testing.T) {
	handler, _ := testAPI(t)
	for _, tc := range []struct {
		method, path string
		status       int
	}{
		{"GET", "/healthz", 200},
		{"GET", "/missing", 404},
		{"GET", "/notes/x/extra", 404},
		{"PATCH", "/notes/x", 405},
		{"DELETE", "/notes", 405},
		{"POST", "/healthz", 405},
	} {
		w := request(handler, tc.method, tc.path, "", "")
		if w.Code != tc.status {
			t.Errorf("%s %s: %d %s", tc.method, tc.path, w.Code, w.Body.String())
		}
		if tc.status == 405 && w.Header().Get("Allow") == "" {
			t.Error("405 response missing Allow header")
		}
	}
}

func TestStorageFailureResponse(t *testing.T) {
	handler, path := testAPI(t)
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(path, 0700); err != nil {
		t.Fatal(err)
	}
	w := request(handler, "POST", "/notes", `{"title":"Cannot save"}`, "application/json")
	if w.Code != http.StatusInternalServerError || strings.Contains(w.Body.String(), path) {
		t.Fatalf("storage error response: %d %s", w.Code, w.Body.String())
	}
}
