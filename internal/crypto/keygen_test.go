package crypto

import (
	"encoding/hex"
	"strings"
	"testing"
)

func TestGenerateKey_Length(t *testing.T) {
	key, err := GenerateKey()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// hex-encoded 32 bytes = 64 characters
	if len(key) != 64 {
		t.Errorf("expected key length 64, got %d", len(key))
	}
}

func TestGenerateKey_IsHex(t *testing.T) {
	key, err := GenerateKey()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := hex.DecodeString(key); err != nil {
		t.Errorf("key is not valid hex: %v", err)
	}
}

func TestGenerateKey_Uniqueness(t *testing.T) {
	a, _ := GenerateKey()
	b, _ := GenerateKey()
	if a == b {
		t.Error("expected two generated keys to differ")
	}
}

func TestDeriveKey_Length(t *testing.T) {
	key := DeriveKey("passphrase", "somesalt")
	if len(key) != KeyLength {
		t.Errorf("expected derived key length %d, got %d", KeyLength, len(key))
	}
}

func TestDeriveKey_Deterministic(t *testing.T) {
	a := DeriveKey("secret", "salt123")
	b := DeriveKey("secret", "salt123")
	if string(a) != string(b) {
		t.Error("DeriveKey should be deterministic for same inputs")
	}
}

func TestDeriveKey_DifferentSalts(t *testing.T) {
	a := DeriveKey("secret", "saltA")
	b := DeriveKey("secret", "saltB")
	if string(a) == string(b) {
		t.Error("different salts should produce different keys")
	}
}

func TestDeriveKey_DifferentPassphrases(t *testing.T) {
	a := DeriveKey("passA", "salt")
	b := DeriveKey("passB", "salt")
	if string(a) == string(b) {
		t.Error("different passphrases should produce different keys")
	}
}

func TestFingerprintKey_Length(t *testing.T) {
	fp := FingerprintKey("somekey")
	// 4 bytes hex-encoded = 8 characters
	if len(fp) != 8 {
		t.Errorf("expected fingerprint length 8, got %d", len(fp))
	}
}

func TestFingerprintKey_IsHex(t *testing.T) {
	fp := FingerprintKey("somekey")
	if _, err := hex.DecodeString(fp); err != nil {
		t.Errorf("fingerprint is not valid hex: %v", err)
	}
}

func TestFingerprintKey_Consistent(t *testing.T) {
	a := FingerprintKey("mykey")
	b := FingerprintKey("mykey")
	if a != b {
		t.Error("fingerprint should be consistent for same input")
	}
}

func TestFingerprintKey_IsLowercase(t *testing.T) {
	fp := FingerprintKey("testkey")
	if fp != strings.ToLower(fp) {
		t.Errorf("expected lowercase fingerprint, got %q", fp)
	}
}
