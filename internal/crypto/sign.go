package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

// ErrInvalidSignature is returned when signature verification fails.
var ErrInvalidSignature = errors.New("invalid signature")

// Signer provides HMAC-based signing and verification of data.
type Signer struct {
	key []byte
}

// NewSigner creates a new Signer using the provided passphrase.
// The passphrase is derived into a fixed-length key via DeriveKey.
func NewSigner(passphrase string) (*Signer, error) {
	if passphrase == "" {
		return nil, errors.New("passphrase must not be empty")
	}
	key, err := DeriveKey(passphrase)
	if err != nil {
		return nil, err
	}
	return &Signer{key: key}, nil
}

// Sign computes an HMAC-SHA256 signature over data and returns it as a hex string.
func (s *Signer) Sign(data []byte) string {
	mac := hmac.New(sha256.New, s.key)
	mac.Write(data)
	return hex.EncodeToString(mac.Sum(nil))
}

// Verify checks that sig is the valid HMAC-SHA256 signature for data.
// Returns ErrInvalidSignature if the signature does not match.
func (s *Signer) Verify(data []byte, sig string) error {
	expected := s.Sign(data)
	if !hmac.Equal([]byte(expected), []byte(sig)) {
		return ErrInvalidSignature
	}
	return nil
}
