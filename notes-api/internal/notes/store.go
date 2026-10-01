package notes

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

var ErrNotFound = errors.New("note not found")

type ValidationError struct{ Message string }

func (e *ValidationError) Error() string { return e.Message }

type Note struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Input struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (in Input) validate() (Input, error) {
	in.Title = strings.TrimSpace(in.Title)
	if in.Title == "" || utf8.RuneCountInString(in.Title) > 200 {
		return Input{}, &ValidationError{"title must contain between 1 and 200 characters"}
	}
	if utf8.RuneCountInString(in.Content) > 10000 {
		return Input{}, &ValidationError{"content must contain at most 10000 characters"}
	}
	return in, nil
}

// Store serializes mutations and persists them before publishing changes to readers.
// A data file must be used by only one server process at a time.
type Store struct {
	mu    sync.RWMutex
	path  string
	notes map[string]Note
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("data file path must not be empty")
	}
	s := &Store{path: path, notes: make(map[string]Note)}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		// Check the location is writable before accepting requests.
		if err := s.persist(s.notes); err != nil {
			return nil, err
		}
		return s, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read notes: %w", err)
	}
	var saved []Note
	if err := json.Unmarshal(data, &saved); err != nil {
		return nil, fmt.Errorf("decode notes: %w", err)
	}
	if saved == nil {
		return nil, errors.New("notes data must be a JSON array")
	}
	for _, note := range saved {
		_, validationErr := (Input{Title: note.Title, Content: note.Content}).validate()
		if note.ID == "" || validationErr != nil || note.CreatedAt.IsZero() || note.UpdatedAt.Before(note.CreatedAt) {
			return nil, errors.New("notes data contains an invalid note")
		}
		if _, exists := s.notes[note.ID]; exists {
			return nil, errors.New("notes data contains duplicate IDs")
		}
		s.notes[note.ID] = note
	}
	return s, nil
}

func ordered(notes map[string]Note) []Note {
	result := make([]Note, 0, len(notes))
	for _, note := range notes {
		result = append(result, note)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].CreatedAt.Equal(result[j].CreatedAt) {
			return result[i].ID < result[j].ID
		}
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})
	return result
}

func (s *Store) List() []Note {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return ordered(s.notes)
}

func (s *Store) Get(id string) (Note, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	note, exists := s.notes[id]
	if !exists {
		return Note{}, ErrNotFound
	}
	return note, nil
}

func (s *Store) Create(input Input) (Note, error) {
	input, err := input.validate()
	if err != nil {
		return Note{}, err
	}
	var id [16]byte
	if _, err := rand.Read(id[:]); err != nil {
		return Note{}, fmt.Errorf("generate note ID: %w", err)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now().UTC()
	note := Note{ID: hex.EncodeToString(id[:]), Title: input.Title, Content: input.Content, CreatedAt: now, UpdatedAt: now}
	if _, exists := s.notes[note.ID]; exists {
		return Note{}, errors.New("generated duplicate note ID")
	}
	return note, s.replace(note.ID, &note)
}

func (s *Store) Update(id string, input Input) (Note, error) {
	input, err := input.validate()
	if err != nil {
		return Note{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	note, exists := s.notes[id]
	if !exists {
		return Note{}, ErrNotFound
	}
	note.Title, note.Content, note.UpdatedAt = input.Title, input.Content, time.Now().UTC()
	return note, s.replace(id, &note)
}

func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.notes[id]; !exists {
		return ErrNotFound
	}
	return s.replace(id, nil)
}

// replace is called with the write lock held. Failed saves leave memory unchanged.
func (s *Store) replace(id string, note *Note) error {
	next := make(map[string]Note, len(s.notes)+1)
	for key, value := range s.notes {
		next[key] = value
	}
	if note == nil {
		delete(next, id)
	} else {
		next[id] = *note
	}
	if err := s.persist(next); err != nil {
		return err
	}
	s.notes = next
	return nil
}

func (s *Store) persist(notes map[string]Note) error {
	data, err := json.MarshalIndent(ordered(notes), "", "  ")
	if err != nil {
		return fmt.Errorf("encode notes: %w", err)
	}
	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("create data directory: %w", err)
	}
	file, err := os.CreateTemp(dir, ".notes-*.tmp")
	if err != nil {
		return fmt.Errorf("create temporary data file: %w", err)
	}
	defer os.Remove(file.Name())
	defer file.Close()
	if _, err := file.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("write notes: %w", err)
	}
	if err := file.Sync(); err != nil {
		return fmt.Errorf("sync notes: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close notes: %w", err)
	}
	if err := os.Rename(file.Name(), s.path); err != nil {
		return fmt.Errorf("replace notes data: %w", err)
	}
	return nil
}
