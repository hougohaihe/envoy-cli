package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

const (
	KeyLength  = 32
	SaltLength = 16
	Iterations = 100_000
)

// GenerateKey generates a cryptographically secure random hex-encoded key.
func GenerateKey() (string, error) {
	buf := make([]byte, KeyLength)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("keygen: failed to generate random key: %w", err)
	}
	return hex.EncodeToString(buf), nil
}

// DeriveKey derives a deterministic key from a passphrase and salt using PBKDF2-SHA256.
func DeriveKey(passphrase, salt string) []byte {
	return pbkdf2.Key(
		[]byte(passphrase),
		[]byte(salt),
		Iterations,
		KeyLength,
		sha256.New,
	)
}

// FingerprintKey returns a short hex fingerprint of a key for identification.
func FingerprintKey(key string) string {
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:4])
}
