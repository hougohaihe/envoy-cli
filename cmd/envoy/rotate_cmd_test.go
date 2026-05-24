package main

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/user/envoy-cli/internal/crypto"
)

func setupRotateVault(t *testing.T, pass string) (string, *crypto.Vault) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "vault.enc")
	v := crypto.NewVault(path)
	if err := v.Save(map[string]string{"HELLO": "world"}, pass); err != nil {
		t.Fatalf("setup vault: %v", err)
	}
	return path, v
}

func TestRotateCmd_Success(t *testing.T) {
	path, v := setupRotateVault(t, "old")

	cmd := buildRotateCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{
		"--old-passphrase", "old",
		"--new-passphrase", "new",
		"--vault", path,
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}

	if got := out.String(); got == "" {
		t.Error("expected success message in stdout")
	}

	// Verify new passphrase works.
	data, err := v.Load("new")
	if err != nil {
		t.Fatalf("load with new pass: %v", err)
	}
	if data["HELLO"] != "world" {
		t.Errorf("unexpected data after rotate: %v", data)
	}
}

func TestRotateCmd_MissingOldPassphrase(t *testing.T) {
	path, _ := setupRotateVault(t, "old")

	cmd := buildRotateCmd()
	cmd.SetArgs([]string{"--new-passphrase", "new", "--vault", path})

	if err := cmd.Execute(); err == nil {
		t.Error("expected error when --old-passphrase is missing")
	}
}

func TestRotateCmd_MissingNewPassphrase(t *testing.T) {
	path, _ := setupRotateVault(t, "old")

	cmd := buildRotateCmd()
	cmd.SetArgs([]string{"--old-passphrase", "old", "--vault", path})

	if err := cmd.Execute(); err == nil {
		t.Error("expected error when --new-passphrase is missing")
	}
}

func TestRotateCmd_WrongOldPassphrase(t *testing.T) {
	path, _ := setupRotateVault(t, "correct")

	cmd := buildRotateCmd()
	cmd.SetArgs([]string{
		"--old-passphrase", "wrong",
		"--new-passphrase", "new",
		"--vault", path,
	})

	if err := cmd.Execute(); err == nil {
		t.Error("expected error with wrong old passphrase")
	}
}
