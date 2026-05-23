package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/envoy-cli/internal/config"
)

func tempConfigPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "config.json")
}

func TestLoadFrom_FileNotExist(t *testing.T) {
	cfg, err := config.LoadFrom("/nonexistent/path/config.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil default config")
	}
	if cfg.Aliases == nil {
		t.Error("expected Aliases map to be initialised")
	}
}

func TestSaveTo_And_LoadFrom(t *testing.T) {
	path := tempConfigPath(t)

	orig := &config.Config{
		RemoteURL: "https://api.example.com",
		ActiveEnv: "staging",
		AuthToken: "tok_abc123",
		Aliases:   map[string]string{"prod": "production"},
	}

	if err := orig.SaveTo(path); err != nil {
		t.Fatalf("SaveTo failed: %v", err)
	}

	loaded, err := config.LoadFrom(path)
	if err != nil {
		t.Fatalf("LoadFrom failed: %v", err)
	}

	if loaded.RemoteURL != orig.RemoteURL {
		t.Errorf("RemoteURL: want %q, got %q", orig.RemoteURL, loaded.RemoteURL)
	}
	if loaded.ActiveEnv != orig.ActiveEnv {
		t.Errorf("ActiveEnv: want %q, got %q", orig.ActiveEnv, loaded.ActiveEnv)
	}
	if loaded.AuthToken != orig.AuthToken {
		t.Errorf("AuthToken: want %q, got %q", orig.AuthToken, loaded.AuthToken)
	}
	if loaded.Aliases["prod"] != "production" {
		t.Errorf("Aliases[prod]: want %q, got %q", "production", loaded.Aliases["prod"])
	}
}

func TestSaveTo_CreatesParentDirs(t *testing.T) {
	base := t.TempDir()
	path := filepath.Join(base, "nested", "dir", "config.json")

	cfg := &config.Config{RemoteURL: "https://example.com", Aliases: map[string]string{}}
	if err := cfg.SaveTo(path); err != nil {
		t.Fatalf("SaveTo failed: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected config file to exist: %v", err)
	}
}

func TestLoadFrom_InvalidJSON(t *testing.T) {
	path := tempConfigPath(t)
	if err := os.WriteFile(path, []byte("not-json"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := config.LoadFrom(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
