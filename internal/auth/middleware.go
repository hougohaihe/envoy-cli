package auth

import (
	"errors"
	"time"
)

// ErrTokenExpired is returned when a token has passed its expiry time.
var ErrTokenExpired = errors.New("auth: token has expired")

// ErrTokenInvalid is returned when a token is malformed or missing.
var ErrTokenInvalid = errors.New("auth: token is invalid")

// Validator checks whether a token is valid and not expired.
type Validator struct {
	manager *Manager
}

// NewValidator creates a Validator backed by the given Manager.
func NewValidator(m *Manager) *Validator {
	return &Validator{manager: m}
}

// Validate loads the stored token and verifies it is still valid.
// Returns the Token on success, or an error if missing or expired.
func (v *Validator) Validate() (*Token, error) {
	tok, err := v.manager.Load()
	if err != nil {
		return nil, ErrTokenInvalid
	}
	if tok == nil {
		return nil, ErrTokenInvalid
	}
	if time.Now().After(tok.ExpiresAt) {
		return nil, ErrTokenExpired
	}
	return tok, nil
}

// IsAuthenticated returns true if a valid, non-expired token is present.
func (v *Validator) IsAuthenticated() bool {
	_, err := v.Validate()
	return err == nil
}

// RequireAuth calls Validate and returns the token or a descriptive error
// suitable for use as a pre-flight check in CLI commands.
func (v *Validator) RequireAuth() (*Token, error) {
	tok, err := v.Validate()
	if errors.Is(err, ErrTokenExpired) {
		return nil, errors.New("session expired — please log in again")
	}
	if errors.Is(err, ErrTokenInvalid) {
		return nil, errors.New("not authenticated — please log in first")
	}
	return tok, nil
}
