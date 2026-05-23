package crypto

import (
	"os"
	"path/filepath"
	"testing"
)

func tempVaultPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "sub", "vault.json")
}

func TestVault_SaveAndLoad(t *testing.T) {
	path := tempVaultPath(t)
	v := NewVault(path, "my-passphrase")

	origData := []byte("API_KEY=secret\nDB_URL=postgres://localhost/test")
	if err := v.Save("production", origData); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	name, data, err := v.Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if name != "production" {
		t.Errorf("Load() name = %q, want %q", name, "production")
	}
	if string(data) != string(origData) {
		t.Errorf("Load() data = %q, want %q", data, origData)
	}
}

func TestVault_Save_CreatesParentDirs(t *testing.T) {
	path := tempVaultPath(t)
	v := NewVault(path, "pass")

	if err := v.Save("dev", []byte("FOO=bar")); err != nil {
		t.Fatalf("Save() should create parent dirs, got error: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("vault file was not created")
	}
}

func TestVault_Load_WrongPassphrase(t *testing.T) {
	path := tempVaultPath(t)

	v1 := NewVault(path, "correct")
	if err := v1.Save("staging", []byte("SECRET=abc")); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	v2 := NewVault(path, "wrong")
	_, _, err := v2.Load()
	if err == nil {
		t.Error("Load() expected error with wrong passphrase, got nil")
	}
}

func TestVault_Load_FileNotExist(t *testing.T) {
	v := NewVault("/nonexistent/path/vault.json", "pass")
	_, _, err := v.Load()
	if err == nil {
		t.Error("Load() expected error for missing file, got nil")
	}
}
