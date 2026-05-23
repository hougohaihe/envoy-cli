package crypto_test

import (
	"encoding/hex"
	"testing"

	"github.com/your-org/envoy-cli/internal/crypto"
)

func TestGenerateKey_Length(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// hex-encoded 32 bytes = 64 chars
	if len(key) != 64 {
		t.Errorf("expected length 64, got %d", len(key))
	}
}

func TestGenerateKey_IsHex(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := hex.DecodeString(key); err != nil {
		t.Errorf("key is not valid hex: %v", err)
	}
}

func TestGenerateKey_Uniqueness(t *testing.T) {
	a, _ := crypto.GenerateKey()
	b, _ := crypto.GenerateKey()
	if a == b {
		t.Error("expected unique keys, got identical values")
	}
}

func TestDeriveKey_Length(t *testing.T) {
	key := crypto.DeriveKey("passphrase", "somesalt")
	if len(key) != 32 {
		t.Errorf("expected 32 bytes, got %d", len(key))
	}
}

func TestDeriveKey_Deterministic(t *testing.T) {
	a := crypto.DeriveKey("secret", "salt123")
	b := crypto.DeriveKey("secret", "salt123")
	if string(a) != string(b) {
		t.Error("DeriveKey should be deterministic for same inputs")
	}
}

func TestDeriveKey_DifferentSalts(t *testing.T) {
	a := crypto.DeriveKey("secret", "salt1")
	b := crypto.DeriveKey("secret", "salt2")
	if string(a) == string(b) {
		t.Error("different salts should produce different keys")
	}
}

func TestDeriveKey_DifferentPassphrases(t *testing.T) {
	a := crypto.DeriveKey("pass1", "salt")
	b := crypto.DeriveKey("pass2", "salt")
	if string(a) == string(b) {
		t.Error("different passphrases should produce different keys")
	}
}

func TestFingerprintKey_Length(t *testing.T) {
	fp := crypto.FingerprintKey("somekey")
	// 4 bytes hex-encoded = 8 chars
	if len(fp) != 8 {
		t.Errorf("expected fingerprint length 8, got %d", len(fp))
	}
}

func TestFingerprintKey_Deterministic(t *testing.T) {
	a := crypto.FingerprintKey("mykey")
	b := crypto.FingerprintKey("mykey")
	if a != b {
		t.Error("fingerprint should be deterministic")
	}
}

func TestFingerprintKey_DifferentKeys(t *testing.T) {
	a := crypto.FingerprintKey("key1")
	b := crypto.FingerprintKey("key2")
	if a == b {
		t.Error("different keys should produce different fingerprints")
	}
}
