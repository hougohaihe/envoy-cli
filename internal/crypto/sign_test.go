package crypto

import (
	"testing"
)

func TestNewSigner_EmptyPassphrase(t *testing.T) {
	_, err := NewSigner("")
	if err == nil {
		t.Fatal("expected error for empty passphrase, got nil")
	}
}

func TestSign_ProducesHexString(t *testing.T) {
	s, err := NewSigner("test-passphrase")
	if err != nil {
		t.Fatalf("NewSigner: %v", err)
	}
	sig := s.Sign([]byte("hello world"))
	if len(sig) == 0 {
		t.Fatal("expected non-empty signature")
	}
	// SHA256 HMAC hex is 64 chars
	if len(sig) != 64 {
		t.Fatalf("expected 64-char hex signature, got %d", len(sig))
	}
}

func TestVerify_ValidSignature(t *testing.T) {
	s, err := NewSigner("my-secret")
	if err != nil {
		t.Fatalf("NewSigner: %v", err)
	}
	data := []byte("payload data")
	sig := s.Sign(data)
	if err := s.Verify(data, sig); err != nil {
		t.Fatalf("expected valid signature, got: %v", err)
	}
}

func TestVerify_InvalidSignature(t *testing.T) {
	s, err := NewSigner("my-secret")
	if err != nil {
		t.Fatalf("NewSigner: %v", err)
	}
	data := []byte("payload data")
	err = s.Verify(data, "deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
	if err == nil {
		t.Fatal("expected ErrInvalidSignature, got nil")
	}
	if err != ErrInvalidSignature {
		t.Fatalf("expected ErrInvalidSignature, got: %v", err)
	}
}

func TestVerify_TamperedData(t *testing.T) {
	s, err := NewSigner("my-secret")
	if err != nil {
		t.Fatalf("NewSigner: %v", err)
	}
	original := []byte("original payload")
	sig := s.Sign(original)
	tampered := []byte("tampered payload")
	if err := s.Verify(tampered, sig); err != ErrInvalidSignature {
		t.Fatalf("expected ErrInvalidSignature for tampered data, got: %v", err)
	}
}

func TestSign_DifferentPassphraseDifferentSig(t *testing.T) {
	s1, _ := NewSigner("passphrase-one")
	s2, _ := NewSigner("passphrase-two")
	data := []byte("same data")
	if s1.Sign(data) == s2.Sign(data) {
		t.Fatal("expected different signatures for different passphrases")
	}
}
