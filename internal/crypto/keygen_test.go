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
	// KeySize bytes → hex string is 2× as long
	expected := KeySize * 2
	if len(key) != expected {
		t.Errorf("expected key length %d, got %d", expected, len(key))
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
		t.Error("two generated keys should not be equal")
	}
}

func TestDeriveKey_Length(t *testing.T) {
	key, err := DeriveKey("mysecretpassphrase")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(key) != 32 {
		t.Errorf("expected 32-byte key, got %d", len(key))
	}
}

func TestDeriveKey_Deterministic(t *testing.T) {
	a, _ := DeriveKey("same-pass")
	b, _ := DeriveKey("same-pass")
	if string(a) != string(b) {
		t.Error("DeriveKey should be deterministic for the same passphrase")
	}
}

func TestDeriveKey_DifferentInputs(t *testing.T) {
	a, _ := DeriveKey("pass-one")
	b, _ := DeriveKey("pass-two")
	if string(a) == string(b) {
		t.Error("different passphrases should produce different keys")
	}
}

func TestDeriveKey_EmptyPassphrase(t *testing.T) {
	_, err := DeriveKey("")
	if err == nil {
		t.Error("expected error for empty passphrase")
	}
}

func TestFingerprintKey_Format(t *testing.T) {
	fp := FingerprintKey("some-key-material")
	// 4 bytes → 8 hex chars
	if len(fp) != 8 {
		t.Errorf("expected fingerprint length 8, got %d", len(fp))
	}
	if strings.ToLower(fp) != fp {
		t.Error("fingerprint should be lowercase hex")
	}
}

func TestFingerprintKey_Deterministic(t *testing.T) {
	a := FingerprintKey("key")
	b := FingerprintKey("key")
	if a != b {
		t.Error("fingerprint should be deterministic")
	}
}
