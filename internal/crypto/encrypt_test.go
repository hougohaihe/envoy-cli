package crypto

import (
	"bytes"
	"testing"
)

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	enc := NewEncryptor("supersecretpassphrase")
	plaintext := []byte("DATABASE_URL=postgres://localhost/dev")

	ciphertext, err := enc.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt() error: %v", err)
	}

	if bytes.Equal(ciphertext, plaintext) {
		t.Fatal("Encrypt() returned plaintext unchanged")
	}

	decrypted, err := enc.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Decrypt() error: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("Decrypt() = %q, want %q", decrypted, plaintext)
	}
}

func TestEncrypt_ProducesUniqueCiphertexts(t *testing.T) {
	enc := NewEncryptor("passphrase")
	plaintext := []byte("API_KEY=abc123")

	c1, _ := enc.Encrypt(plaintext)
	c2, _ := enc.Encrypt(plaintext)

	if bytes.Equal(c1, c2) {
		t.Error("Encrypt() produced identical ciphertexts for same plaintext (nonce not random)")
	}
}

func TestDecrypt_WrongPassphrase(t *testing.T) {
	enc := NewEncryptor("correct-passphrase")
	ciphertext, _ := enc.Encrypt([]byte("SECRET=value"))

	wrong := NewEncryptor("wrong-passphrase")
	_, err := wrong.Decrypt(ciphertext)
	if err == nil {
		t.Error("Decrypt() expected error with wrong passphrase, got nil")
	}
}

func TestDecrypt_TooShort(t *testing.T) {
	enc := NewEncryptor("passphrase")
	_, err := enc.Decrypt([]byte{0x01, 0x02})
	if err == nil {
		t.Error("Decrypt() expected error for short ciphertext, got nil")
	}
}

func TestNewEncryptor_DifferentPassphrasesDifferentKeys(t *testing.T) {
	e1 := NewEncryptor("pass1")
	e2 := NewEncryptor("pass2")

	plaintext := []byte("VAR=value")
	c1, _ := e1.Encrypt(plaintext)

	_, err := e2.Decrypt(c1)
	if err == nil {
		t.Error("expected decryption to fail with different passphrase")
	}
}
