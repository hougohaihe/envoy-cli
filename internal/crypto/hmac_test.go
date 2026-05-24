package crypto

import (
	"testing"
)

func TestNewHMACSigner_EmptyPassphrase(t *testing.T) {
	_, err := NewHMACSigner("")
	if err == nil {
		t.Fatal("expected error for empty passphrase")
	}
}

func TestHMAC_SignProducesHexString(t *testing.T) {
	s, err := NewHMACSigner("secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	mac := s.Sign([]byte("hello world"))
	if len(mac) == 0 {
		t.Fatal("expected non-empty MAC")
	}
	// SHA256 HMAC = 32 bytes = 64 hex chars
	if len(mac) != 64 {
		t.Fatalf("expected 64 hex chars, got %d", len(mac))
	}
}

func TestHMAC_VerifyValidSignature(t *testing.T) {
	s, _ := NewHMACSigner("secret")
	data := []byte("envoy payload")
	mac := s.Sign(data)
	if err := s.Verify(data, mac); err != nil {
		t.Fatalf("expected valid signature, got: %v", err)
	}
}

func TestHMAC_VerifyInvalidSignature(t *testing.T) {
	s, _ := NewHMACSigner("secret")
	data := []byte("envoy payload")
	if err := s.Verify(data, "deadbeef"); err == nil {
		t.Fatal("expected error for invalid MAC")
	}
}

func TestHMAC_VerifyTamperedData(t *testing.T) {
	s, _ := NewHMACSigner("secret")
	original := []byte("original")
	mac := s.Sign(original)
	if err := s.Verify([]byte("tampered"), mac); err == nil {
		t.Fatal("expected error for tampered data")
	}
}

func TestHMAC_DifferentPassphrasesProduceDifferentMACs(t *testing.T) {
	s1, _ := NewHMACSigner("passphrase-one")
	s2, _ := NewHMACSigner("passphrase-two")
	data := []byte("same data")
	if s1.Sign(data) == s2.Sign(data) {
		t.Fatal("different passphrases should produce different MACs")
	}
}

func TestHMAC_Deterministic(t *testing.T) {
	s, _ := NewHMACSigner("deterministic-key")
	data := []byte("consistent")
	if s.Sign(data) != s.Sign(data) {
		t.Fatal("HMAC should be deterministic for the same key and data")
	}
}
