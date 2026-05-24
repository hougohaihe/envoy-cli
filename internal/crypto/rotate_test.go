package crypto_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envoy-cli/internal/crypto"
)

func tempRotatePath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "vault.enc")
}

func TestRotate_Success(t *testing.T) {
	path := tempRotatePath(t)
	v := crypto.NewVault(path)

	original := map[string]string{"KEY": "value", "FOO": "bar"}
	if err := v.Save(original, "old-pass"); err != nil {
		t.Fatalf("save: %v", err)
	}

	r := crypto.NewRotator(v)
	if err := r.Rotate("old-pass", "new-pass"); err != nil {
		t.Fatalf("rotate: %v", err)
	}

	// Must load with new passphrase.
	loaded, err := v.Load("new-pass")
	if err != nil {
		t.Fatalf("load after rotate: %v", err)
	}
	if loaded["KEY"] != "value" || loaded["FOO"] != "bar" {
		t.Errorf("data mismatch after rotate: %v", loaded)
	}
}

func TestRotate_OldPassphraseNoLongerWorks(t *testing.T) {
	path := tempRotatePath(t)
	v := crypto.NewVault(path)

	if err := v.Save(map[string]string{"X": "1"}, "old-pass"); err != nil {
		t.Fatalf("save: %v", err)
	}

	r := crypto.NewRotator(v)
	_ = r.Rotate("old-pass", "new-pass")

	_, err := v.Load("old-pass")
	if err == nil {
		t.Error("expected error loading with old passphrase after rotation")
	}
}

func TestRotate_SamePassphrase(t *testing.T) {
	path := tempRotatePath(t)
	v := crypto.NewVault(path)
	_ = v.Save(map[string]string{}, "pass")

	r := crypto.NewRotator(v)
	err := r.Rotate("pass", "pass")
	if err == nil {
		t.Error("expected error when old and new passphrases are equal")
	}
}

func TestRotate_WrongOldPassphrase(t *testing.T) {
	path := tempRotatePath(t)
	v := crypto.NewVault(path)
	_ = v.Save(map[string]string{"A": "B"}, "correct")

	r := crypto.NewRotator(v)
	err := r.Rotate("wrong", "new-pass")
	if err == nil {
		t.Error("expected error with wrong old passphrase")
	}
}

func TestRotate_EmptyPassphrase(t *testing.T) {
	path := tempRotatePath(t)
	v := crypto.NewVault(path)
	_ = v.Save(map[string]string{}, "pass")

	r := crypto.NewRotator(v)
	if err := r.Rotate("", "new"); err == nil {
		t.Error("expected error for empty old passphrase")
	}
	if err := r.Rotate("pass", ""); err == nil {
		t.Error("expected error for empty new passphrase")
	}
}

func TestRotate_FileNotExist(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent.enc")
	v := crypto.NewVault(path)
	r := crypto.NewRotator(v)
	err := r.Rotate("old", "new")
	if err == nil {
		t.Error("expected error when vault file does not exist")
	}
	if _, statErr := os.Stat(path); !os.IsNotExist(statErr) {
		t.Error("vault file should not have been created on failed rotation")
	}
}
