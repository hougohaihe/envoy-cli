package auth

import (
	"testing"
	"time"
)

func newTestValidator(t *testing.T) (*Validator, *Manager) {
	t.Helper()
	path := tempStorePath(t)
	m, err := NewManagerWithStore(NewTokenStoreAt(path))
	if err != nil {
		t.Fatalf("NewManagerWithStore: %v", err)
	}
	return NewValidator(m), m
}

func TestValidate_NoToken(t *testing.T) {
	v, _ := newTestValidator(t)
	_, err := v.Validate()
	if err != ErrTokenInvalid {
		t.Fatalf("expected ErrTokenInvalid, got %v", err)
	}
}

func TestValidate_ValidToken(t *testing.T) {
	v, m := newTestValidator(t)
	tok, _ := NewToken("user@example.com", time.Hour)
	if err := m.Save(tok); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := v.Validate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Identity != tok.Identity {
		t.Errorf("identity mismatch: got %q want %q", got.Identity, tok.Identity)
	}
}

func TestValidate_ExpiredToken(t *testing.T) {
	v, m := newTestValidator(t)
	tok, _ := NewToken("user@example.com", time.Millisecond)
	if err := m.Save(tok); err != nil {
		t.Fatalf("Save: %v", err)
	}
	time.Sleep(5 * time.Millisecond)
	_, err := v.Validate()
	if err != ErrTokenExpired {
		t.Fatalf("expected ErrTokenExpired, got %v", err)
	}
}

func TestIsAuthenticated_True(t *testing.T) {
	v, m := newTestValidator(t)
	tok, _ := NewToken("user@example.com", time.Hour)
	_ = m.Save(tok)
	if !v.IsAuthenticated() {
		t.Error("expected IsAuthenticated to return true")
	}
}

func TestIsAuthenticated_False(t *testing.T) {
	v, _ := newTestValidator(t)
	if v.IsAuthenticated() {
		t.Error("expected IsAuthenticated to return false when no token stored")
	}
}

func TestRequireAuth_NotLoggedIn(t *testing.T) {
	v, _ := newTestValidator(t)
	_, err := v.RequireAuth()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "not authenticated — please log in first" {
		t.Errorf("unexpected message: %q", err.Error())
	}
}

func TestRequireAuth_Expired(t *testing.T) {
	v, m := newTestValidator(t)
	tok, _ := NewToken("user@example.com", time.Millisecond)
	_ = m.Save(tok)
	time.Sleep(5 * time.Millisecond)
	_, err := v.RequireAuth()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "session expired — please log in again" {
		t.Errorf("unexpected message: %q", err.Error())
	}
}
