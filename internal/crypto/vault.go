package crypto

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Vault persists encrypted env set payloads to disk.
type Vault struct {
	path      string
	encryptor *Encryptor
}

type vaultEntry struct {
	Name    string `json:"name"`
	Payload string `json:"payload"` // base64-encoded ciphertext
}

// NewVault creates a Vault that stores data at path, encrypted with passphrase.
func NewVault(path, passphrase string) *Vault {
	return &Vault{
		path:      path,
		encryptor: NewEncryptor(passphrase),
	}
}

// Save encrypts data and writes it to the vault file for the given name.
func (v *Vault) Save(name string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(v.path), 0700); err != nil {
		return fmt.Errorf("vault: mkdir: %w", err)
	}

	ciphertext, err := v.encryptor.Encrypt(data)
	if err != nil {
		return fmt.Errorf("vault: encrypt: %w", err)
	}

	entry := vaultEntry{
		Name:    name,
		Payload: base64.StdEncoding.EncodeToString(ciphertext),
	}

	b, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("vault: marshal: %w", err)
	}

	return os.WriteFile(v.path, b, 0600)
}

// Load reads and decrypts the vault file, returning the plaintext data.
func (v *Vault) Load() (string, []byte, error) {
	b, err := os.ReadFile(v.path)
	if err != nil {
		return "", nil, fmt.Errorf("vault: read: %w", err)
	}

	var entry vaultEntry
	if err := json.Unmarshal(b, &entry); err != nil {
		return "", nil, fmt.Errorf("vault: unmarshal: %w", err)
	}

	ciphertext, err := base64.StdEncoding.DecodeString(entry.Payload)
	if err != nil {
		return "", nil, fmt.Errorf("vault: base64 decode: %w", err)
	}

	plaintext, err := v.encryptor.Decrypt(ciphertext)
	if err != nil {
		return "", nil, err
	}

	return entry.Name, plaintext, nil
}
