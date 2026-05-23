package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const defaultTokenFile = ".envoy/token.json"

// TokenStore persists and retrieves auth tokens from disk.
type TokenStore struct {
	path string
}

// NewTokenStore creates a TokenStore at the default path.
func NewTokenStore() *TokenStore {
	home, _ := os.UserHomeDir()
	return &TokenStore{path: filepath.Join(home, defaultTokenFile)}
}

// NewTokenStoreAt creates a TokenStore at a custom path.
func NewTokenStoreAt(path string) *TokenStore {
	return &TokenStore{path: path}
}

// Save writes the token to disk, creating parent directories as needed.
func (s *TokenStore) Save(t *Token) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0700); err != nil {
		return fmt.Errorf("creating token dir: %w", err)
	}
	data, err := json.Marshal(t)
	if err != nil {
		return fmt.Errorf("marshalling token: %w", err)
	}
	return os.WriteFile(s.path, data, 0600)
}

// Load reads the token from disk. Returns ErrNoToken if not found.
func (s *TokenStore) Load() (*Token, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrNoToken
	}
	if err != nil {
		return nil, fmt.Errorf("reading token file: %w", err)
	}
	var t Token
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, fmt.Errorf("parsing token file: %w", err)
	}
	return &t, nil
}

// Clear removes the token file from disk.
func (s *TokenStore) Clear() error {
	err := os.Remove(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

// ErrNoToken is returned when no token file exists.
var ErrNoToken = errors.New("no token found; please run 'envoy login'")
