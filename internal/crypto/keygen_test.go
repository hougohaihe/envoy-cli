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
	// hex-encoded 32 bytes => 64 chars
	if len(key) != 64 {
		t.Errorf("expected key length 64, got %d", len(key))
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
		t.Error("expected two generated keys to differ")
	}
}

func TestDeriveKey_Length(t *testing.T) {
	key := crypto.DeriveKey("passphrase", "somesalt")
	if len(key) != crypto.KeyLength {
		t.Errorf("expected derived key length %d, got %d", crypto.KeyLength, len(key))
	}
}

func TestDeriveKey_Deterministic(t *testing.T) {
	a := crypto.DeriveKey("secret", "salt123")
	b := crypto.DeriveKey("secret", "salt123")
	if string(a) != string(b) {
		t.Error("expected DeriveKey to be deterministic for same inputs")
	}
}

func TestDeriveKey_DifferentSalts(t *testing.T) {
	a := crypto.DeriveKey("secret", "saltA")
	b := crypto.DeriveKey("secret", "saltB")
	if string(a) == string(b) {
		t.Error("expected different salts to produce different keys")
	}
}

func TestDeriveKey_DifferentPassphrases(t *testing.T) {
	a := crypto.DeriveKey("passA", "salt")
	b := crypto.DeriveKey("passB", "salt")
	if string(a) == string(b) {
		t.Error("expected different passphrases to produce different keys")
	}
}

func TestFingerprintKey_Length(t *testing.T) {
	fp := crypto.FingerprintKey("somekey")
	// 8 bytes hex-encoded => 16 chars
	if len(fp) != 16 {
		t.Errorf("expected fingerprint length 16, got %d", len(fp))
	}
}

func TestFingerprintKey_Deterministic(t *testing.T) {
	a := crypto.FingerprintKey("mykey")
	b := crypto.FingerprintKey("mykey")
	if a != b {
		t.Error("expected FingerprintKey to be deterministic")
	}
}

func TestFingerprintKey_Distinct(t *testing.T) {
	a := crypto.FingerprintKey("keyA")
	b := crypto.FingerprintKey("keyB")
	if a == b {
		t.Error("expected different keys to produce different fingerprints")
	}
}
