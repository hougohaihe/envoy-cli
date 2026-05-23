package auth

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempStorePath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "token.json")
}

func TestNewToken_Valid(t *testing.T) {
	tok, err := NewToken(1 * time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Value == "" {
		t.Error("expected non-empty token value")
	}
	if tok.IsExpired() {
		t.Error("new token should not be expired")
	}
	if !tok.IsValid() {
		t.Error("new token should be valid")
	}
}

func TestNewToken_ZeroTTL(t *testing.T) {
	_, err := NewToken(0)
	if err == nil {
		t.Error("expected error for zero TTL")
	}
}

func TestTokenStore_SaveAndLoad(t *testing.T) {
	store := NewTokenStoreAt(tempStorePath(t))
	tok, _ := NewToken(1 * time.Hour)

	if err := store.Save(tok); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if loaded.Value != tok.Value {
		t.Errorf("expected %q, got %q", tok.Value, loaded.Value)
	}
}

func TestTokenStore_LoadMissing(t *testing.T) {
	store := NewTokenStoreAt(filepath.Join(t.TempDir(), "nonexistent.json"))
	_, err := store.Load()
	if err != ErrNoToken {
		t.Errorf("expected ErrNoToken, got %v", err)
	}
}

func TestTokenStore_Clear(t *testing.T) {
	path := tempStorePath(t)
	store := NewTokenStoreAt(path)
	tok, _ := NewToken(1 * time.Hour)
	_ = store.Save(tok)

	if err := store.Clear(); err != nil {
		t.Fatalf("Clear failed: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected token file to be deleted")
	}
}

func TestManager_LoginAndCurrentToken(t *testing.T) {
	store := NewTokenStoreAt(tempStorePath(t))
	mgr := NewManagerWithStore(store)

	tok, err := mgr.Login()
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	current, err := mgr.CurrentToken()
	if err != nil {
		t.Fatalf("CurrentToken failed: %v", err)
	}
	if current.Value != tok.Value {
		t.Errorf("token mismatch: want %q got %q", tok.Value, current.Value)
	}
}

func TestManager_Logout(t *testing.T) {
	store := NewTokenStoreAt(tempStorePath(t))
	mgr := NewManagerWithStore(store)
	_, _ = mgr.Login()

	if err := mgr.Logout(); err != nil {
		t.Fatalf("Logout failed: %v", err)
	}
	_, err := mgr.CurrentToken()
	if err == nil {
		t.Error("expected error after logout")
	}
}
