package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// Token represents an authentication token with metadata.
type Token struct {
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// IsExpired returns true if the token is past its expiry time.
func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// IsValid returns true if the token is non-empty and not expired.
func (t *Token) IsValid() bool {
	return t.Value != "" && !t.IsExpired()
}

// NewToken generates a new random token with the given TTL duration.
func NewToken(ttl time.Duration) (*Token, error) {
	if ttl <= 0 {
		return nil, errors.New("ttl must be positive")
	}
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("generating token: %w", err)
	}
	now := time.Now().UTC()
	return &Token{
		Value:     hex.EncodeToString(b),
		CreatedAt: now,
		ExpiresAt: now.Add(ttl),
	}, nil
}
