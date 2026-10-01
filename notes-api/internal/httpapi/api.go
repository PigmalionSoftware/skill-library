package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"mime"
	"net/http"

	"notes-api/internal/notes"
)

const maxBodyBytes = 64 * 1024

type API struct{ store *notes.Store }

func New(store *notes.Store) http.Handler {
	api := &API{store: store}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /notes", api.list)
	mux.HandleFunc("POST /notes", api.create)
	mux.HandleFunc("GET /notes/{id}", api.get)
	mux.HandleFunc("PUT /notes/{id}", api.update)
	mux.HandleFunc("DELETE /notes/{id}", api.delete)
	mux.HandleFunc("/healthz", methodNotAllowed("GET, HEAD"))
	mux.HandleFunc("/notes", methodNotAllowed("GET, HEAD, POST"))
	mux.HandleFunc("/notes/{id}", methodNotAllowed("GET, HEAD, PUT, DELETE"))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, "endpoint not found")
	})
	return mux
}

func methodNotAllowed(allow string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Allow", allow)
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		slog.Error("write response", "error", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func readInput(w http.ResponseWriter, r *http.Request) (notes.Input, bool) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		writeError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return notes.Input{}, false
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var input notes.Input
	err = decoder.Decode(&input)
	if err == nil {
		var extra any
		err = decoder.Decode(&extra)
		if err == io.EOF {
			return input, true
		}
	}
	var tooLarge *http.MaxBytesError
	if errors.As(err, &tooLarge) {
		writeError(w, http.StatusRequestEntityTooLarge, "request body must be at most 64 KiB")
	} else {
		writeError(w, http.StatusBadRequest, "body must contain one JSON object with only title and content fields")
	}
	return notes.Input{}, false
}

func handleStoreError(w http.ResponseWriter, err error) {
	var validation *notes.ValidationError
	switch {
	case errors.Is(err, notes.ErrNotFound):
		writeError(w, http.StatusNotFound, "note not found")
	case errors.As(err, &validation):
		writeError(w, http.StatusBadRequest, validation.Message)
	default:
		slog.Error("notes storage operation failed", "error", err)
		writeError(w, http.StatusInternalServerError, "could not save notes")
	}
}

func (api *API) list(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, api.store.List())
}

func (api *API) get(w http.ResponseWriter, r *http.Request) {
	note, err := api.store.Get(r.PathValue("id"))
	if err != nil {
		handleStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, note)
}

func (api *API) create(w http.ResponseWriter, r *http.Request) {
	input, ok := readInput(w, r)
	if !ok {
		return
	}
	note, err := api.store.Create(input)
	if err != nil {
		handleStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/notes/"+note.ID)
	writeJSON(w, http.StatusCreated, note)
}

func (api *API) update(w http.ResponseWriter, r *http.Request) {
	input, ok := readInput(w, r)
	if !ok {
		return
	}
	note, err := api.store.Update(r.PathValue("id"), input)
	if err != nil {
		handleStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, note)
}

func (api *API) delete(w http.ResponseWriter, r *http.Request) {
	if err := api.store.Delete(r.PathValue("id")); err != nil {
		handleStoreError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
