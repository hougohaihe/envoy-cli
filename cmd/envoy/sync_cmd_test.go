package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envoy-cli/internal/config"
)

func setupSyncConfig(t *testing.T, remoteURL string) string {
	t.Helper()
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.json")
	cfg := config.Config{
		RemoteURL: remoteURL,
		AuthToken: "test-token",
		VaultPath: filepath.Join(dir, "vault.enc"),
	}
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal config: %v", err)
	}
	if err := os.WriteFile(cfgPath, data, 0600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return cfgPath
}

func TestSyncPushCmd_MissingName(t *testing.T) {
	cmd := buildSyncCmd()
	cmd.SetArgs([]string{"push"})
	var buf bytes.Buffer
	cmd.SetErr(&buf)
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when --name flag is missing")
	}
}

func TestSyncPullCmd_MissingName(t *testing.T) {
	cmd := buildSyncCmd()
	cmd.SetArgs([]string{"pull"})
	var buf bytes.Buffer
	cmd.SetErr(&buf)
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when --name flag is missing")
	}
}

func TestSyncPushCmd_RemoteError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	t.Setenv("ENVOY_CONFIG", setupSyncConfig(t, server.URL))
	t.Setenv("ENVOY_PASSPHRASE", "test-pass")

	cmd := buildRoot()
	cmd.SetArgs([]string{"sync", "push", "--name", "myapp"})
	var buf bytes.Buffer
	cmd.SetErr(&buf)
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error on remote 500")
	}
}

func TestSyncPullCmd_RemoteNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	t.Setenv("ENVOY_CONFIG", setupSyncConfig(t, server.URL))
	t.Setenv("ENVOY_PASSPHRASE", "test-pass")

	cmd := buildRoot()
	cmd.SetArgs([]string{"sync", "pull", "--name", "myapp"})
	var buf bytes.Buffer
	cmd.SetErr(&buf)
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error on remote 404")
	}
}
