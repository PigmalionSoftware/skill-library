package notes

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func openTestStore(t *testing.T) (*Store, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "data", "notes.json")
	store, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	return store, path
}

func TestPersistenceAndLifecycle(t *testing.T) {
	store, path := openTestStore(t)
	if got := store.List(); got == nil || len(got) != 0 {
		t.Fatalf("empty list = %#v", got)
	}
	first, err := store.Create(Input{Title: "  First  ", Content: "Keep whitespace\n"})
	if err != nil {
		t.Fatal(err)
	}
	if first.Title != "First" || len(first.ID) != 32 || first.CreatedAt.IsZero() || first.CreatedAt != first.UpdatedAt {
		t.Fatalf("invalid created note: %+v", first)
	}
	second, err := store.Create(Input{Title: "Second"})
	if err != nil {
		t.Fatal(err)
	}
	listed := store.List()
	if len(listed) != 2 || listed[0].ID != second.ID {
		t.Fatalf("unexpected order: %+v", listed)
	}
	listed[0].Title = "mutated by caller"
	updated, err := store.Update(first.ID, Input{Title: "Edited", Content: "New content"})
	if err != nil {
		t.Fatal(err)
	}
	if updated.CreatedAt != first.CreatedAt || updated.UpdatedAt.Before(first.UpdatedAt) {
		t.Fatalf("invalid update timestamps: %+v", updated)
	}
	reopened, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	got, err := reopened.Get(first.ID)
	if err != nil || got != updated {
		t.Fatalf("reloaded note = %+v, error = %v", got, err)
	}
	got, _ = reopened.Get(second.ID)
	if got.Title != "Second" {
		t.Fatal("list exposed mutable store state")
	}
	if err := reopened.Delete(first.ID); err != nil {
		t.Fatal(err)
	}
	reopened, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := reopened.Get(first.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("deleted note error = %v", err)
	}
	if err := reopened.Delete(first.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("repeated delete error = %v", err)
	}
	if _, err := reopened.Update(first.ID, Input{Title: "missing"}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing update error = %v", err)
	}
}

func TestValidation(t *testing.T) {
	store, _ := openTestStore(t)
	for _, input := range []Input{
		{}, {Title: " \n\t"}, {Title: strings.Repeat("a", 201)},
		{Title: "Valid", Content: strings.Repeat("界", 10001)},
	} {
		if _, err := store.Create(input); err == nil {
			t.Errorf("accepted invalid input %+v", input)
		}
	}
	if len(store.List()) != 0 {
		t.Fatal("invalid input modified the store")
	}
	note, err := store.Create(Input{Title: strings.Repeat("界", 200), Content: strings.Repeat("界", 10000)})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.Update(note.ID, Input{}); err == nil {
		t.Fatal("accepted invalid update")
	}
	got, _ := store.Get(note.ID)
	if got != note {
		t.Fatal("invalid update modified the note")
	}
}

func TestFailedWritesLeaveStateUnchanged(t *testing.T) {
	store, path := openTestStore(t)
	note, err := store.Create(Input{Title: "Original"})
	if err != nil {
		t.Fatal(err)
	}
	// A directory at the destination reliably makes rename fail, even as root.
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(path, 0700); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Create(Input{Title: "Not saved"}); err == nil {
		t.Fatal("create succeeded with unwritable destination")
	}
	if _, err := store.Update(note.ID, Input{Title: "Not saved"}); err == nil {
		t.Fatal("update succeeded with unwritable destination")
	}
	if err := store.Delete(note.ID); err == nil {
		t.Fatal("delete succeeded with unwritable destination")
	}
	got, err := store.Get(note.ID)
	if err != nil || got != note || len(store.List()) != 1 {
		t.Fatalf("failed writes changed state: %+v, %v", got, err)
	}
	leftovers, err := filepath.Glob(filepath.Join(filepath.Dir(path), ".notes-*.tmp"))
	if err != nil || len(leftovers) != 0 {
		t.Fatalf("temporary files leaked: %v, %v", leftovers, err)
	}
}

func TestOpenRejectsCorruptData(t *testing.T) {
	for _, data := range []string{
		"", "broken", "null", "{}", `[{"id":"x","title":"missing timestamps"}]`,
		`[{"id":"x","title":"valid","created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"},{"id":"x","title":"duplicate","created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}]`,
	} {
		path := filepath.Join(t.TempDir(), "notes.json")
		if err := os.WriteFile(path, []byte(data), 0600); err != nil {
			t.Fatal(err)
		}
		if _, err := Open(path); err == nil {
			t.Errorf("accepted corrupt data %q", data)
		}
		saved, err := os.ReadFile(path)
		if err != nil || string(saved) != data {
			t.Fatal("opening corrupt data overwrote the original file")
		}
	}
	if _, err := Open(""); err == nil {
		t.Fatal("accepted empty data path")
	}
}

func TestConcurrentMutations(t *testing.T) {
	store, path := openTestStore(t)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			note, err := store.Create(Input{Title: fmt.Sprintf("Note %d", i)})
			if err != nil {
				t.Error(err)
				return
			}
			if _, err := store.Update(note.ID, Input{Title: "Updated"}); err != nil {
				t.Error(err)
			}
			store.List()
			if _, err := store.Get(note.ID); err != nil {
				t.Error(err)
			}
			if i%2 == 0 {
				if err := store.Delete(note.ID); err != nil {
					t.Error(err)
				}
			}
		}(i)
	}
	wg.Wait()
	reopened, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := reopened.List(); len(got) != 10 {
		t.Fatalf("lost concurrent writes: %d notes", len(got))
	}
}
