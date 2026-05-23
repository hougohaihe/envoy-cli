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
	k1, _ := crypto.GenerateKey()
	k2, _ := crypto.GenerateKey()
	if k1 == k2 {
		t.Error("two generated keys should not be equal")
	}
}

func TestDeriveKey_Length(t *testing.T) {
	key := crypto.DeriveKey("passphrase", "somesalt")
	if len(key) != crypto.KeyLength {
		t.Errorf("expected derived key length %d, got %d", crypto.KeyLength, len(key))
	}
}

func TestDeriveKey_Deterministic(t *testing.T) {
	k1 := crypto.DeriveKey("secret", "salt123")
	k2 := crypto.DeriveKey("secret", "salt123")
	if string(k1) != string(k2) {
		t.Error("DeriveKey should be deterministic for same inputs")
	}
}

func TestDeriveKey_DifferentSalts(t *testing.T) {
	k1 := crypto.DeriveKey("secret", "salt1")
	k2 := crypto.DeriveKey("secret", "salt2")
	if string(k1) == string(k2) {
		t.Error("different salts should produce different keys")
	}
}

func TestDeriveKey_DifferentPassphrases(t *testing.T) {
	k1 := crypto.DeriveKey("pass1", "salt")
	k2 := crypto.DeriveKey("pass2", "salt")
	if string(k1) == string(k2) {
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
	fp1 := crypto.FingerprintKey("mykey")
	fp2 := crypto.FingerprintKey("mykey")
	if fp1 != fp2 {
		t.Error("fingerprint should be deterministic")
	}
}

func TestFingerprintKey_Unique(t *testing.T) {
	fp1 := crypto.FingerprintKey("key1")
	fp2 := crypto.FingerprintKey("key2")
	if fp1 == fp2 {
		t.Error("different keys should have different fingerprints")
	}
}
