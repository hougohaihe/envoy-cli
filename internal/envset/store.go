package envset

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Store persists and retrieves EnvSets from a JSON file on disk.
type Store struct {
	path string
	sets map[string]*EnvSet
}

// NewStore creates a Store backed by the given file path.
func NewStore(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("store path must not be empty")
	}
	s := &Store{path: path, sets: make(map[string]*EnvSet)}
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return s, nil
}

// Save persists an EnvSet into the store and writes it to disk.
func (s *Store) Save(es *EnvSet) error {
	s.sets[es.Name] = es
	return s.flush()
}

// Load retrieves an EnvSet by name from the store.
func (s *Store) Load(name string) (*EnvSet, error) {
	es, ok := s.sets[name]
	if !ok {
		return nil, errors.New("envset not found: " + name)
	}
	return es, nil
}

// List returns the names of all stored EnvSets.
func (s *Store) List() []string {
	names := make([]string, 0, len(s.sets))
	for name := range s.sets {
		names = append(names, name)
	}
	return names
}

// Delete removes an EnvSet from the store and writes changes to disk.
func (s *Store) Delete(name string) error {
	if _, ok := s.sets[name]; !ok {
		return errors.New("envset not found: " + name)
	}
	delete(s.sets, name)
	return s.flush()
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.sets)
}

func (s *Store) flush() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s.sets, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o600)
}
