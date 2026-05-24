package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

// HMACSigner provides HMAC-SHA256 signing and verification.
type HMACSigner struct {
	key []byte
}

// NewHMACSigner creates a new HMACSigner derived from the given passphrase.
// The key is derived using DeriveKey so it is consistent across calls.
func NewHMACSigner(passphrase string) (*HMACSigner, error) {
	if passphrase == "" {
		return nil, errors.New("hmac: passphrase must not be empty")
	}
	key, err := DeriveKey(passphrase)
	if err != nil {
		return nil, err
	}
	return &HMACSigner{key: key}, nil
}

// Sign computes an HMAC-SHA256 over data and returns the hex-encoded MAC.
func (h *HMACSigner) Sign(data []byte) string {
	mac := hmac.New(sha256.New, h.key)
	mac.Write(data)
	return hex.EncodeToString(mac.Sum(nil))
}

// Verify checks that the hex-encoded mac is valid for data.
func (h *HMACSigner) Verify(data []byte, mac string) error {
	expected := h.Sign(data)
	if !hmac.Equal([]byte(expected), []byte(mac)) {
		return errors.New("hmac: signature mismatch")
	}
	return nil
}
