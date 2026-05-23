package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
)

const (
	// KeySize is the number of random bytes used for key generation.
	KeySize = 32
)

// GenerateKey produces a cryptographically secure random hex-encoded key
// suitable for use as a passphrase or shared secret.
func GenerateKey() (string, error) {
	buf := make([]byte, KeySize)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return "", errors.New("keygen: failed to read random bytes: " + err.Error())
	}
	return hex.EncodeToString(buf), nil
}

// DeriveKey produces a deterministic 32-byte key from an arbitrary string
// using SHA-256. Useful for converting user-supplied passphrases into
// fixed-length keys without storing the passphrase itself.
func DeriveKey(passphrase string) ([]byte, error) {
	if passphrase == "" {
		return nil, errors.New("keygen: passphrase must not be empty")
	}
	sum := sha256.Sum256([]byte(passphrase))
	return sum[:], nil
}

// FingerprintKey returns a short hex prefix of the SHA-256 hash of the given
// key material, useful for display/logging without exposing the full secret.
func FingerprintKey(key string) string {
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:4])
}
