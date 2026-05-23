package auth

import (
	"errors"
	"fmt"
	"time"
)

const defaultTokenTTL = 24 * time.Hour

// Manager handles login, logout, and token validation.
type Manager struct {
	store *TokenStore
}

// NewManager returns an auth Manager using the default token store.
func NewManager() *Manager {
	return &Manager{store: NewTokenStore()}
}

// NewManagerWithStore returns an auth Manager with a custom store.
func NewManagerWithStore(store *TokenStore) *Manager {
	return &Manager{store: store}
}

// Login generates a new token and persists it to the store.
func (m *Manager) Login() (*Token, error) {
	t, err := NewToken(defaultTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("creating token: %w", err)
	}
	if err := m.store.Save(t); err != nil {
		return nil, fmt.Errorf("saving token: %w", err)
	}
	return t, nil
}

// Logout removes the stored token.
func (m *Manager) Logout() error {
	return m.store.Clear()
}

// CurrentToken returns the stored token if it exists and is valid.
func (m *Manager) CurrentToken() (*Token, error) {
	t, err := m.store.Load()
	if errors.Is(err, ErrNoToken) {
		return nil, ErrNoToken
	}
	if err != nil {
		return nil, err
	}
	if !t.IsValid() {
		return nil, fmt.Errorf("token expired; please run 'envoy login' again")
	}
	return t, nil
}
