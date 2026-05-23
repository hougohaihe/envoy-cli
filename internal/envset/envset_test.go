package envset

import (
	"testing"
)

func TestNew(t *testing.T) {
	es, err := New("production")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if es.Name != "production" {
		t.Errorf("expected name 'production', got '%s'", es.Name)
	}
	if es.Variables == nil {
		t.Error("expected Variables map to be initialised")
	}
}

func TestNew_EmptyName(t *testing.T) {
	_, err := New("")
	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
}

func TestSet_ValidKey(t *testing.T) {
	es, _ := New("dev")
	if err := es.Set("DATABASE_URL", "postgres://localhost/dev"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	val, ok := es.Get("DATABASE_URL")
	if !ok || val != "postgres://localhost/dev" {
		t.Errorf("expected value 'postgres://localhost/dev', got '%s'" , val)
	}
}

func TestSet_InvalidKey(t *testing.T) {
	es, _ := New("dev")
	if err := es.Set("invalid-key", "value"); err == nil {
		t.Fatal("expected error for invalid key, got nil")
	}
}

func TestDelete(t *testing.T) {
	es, _ := New("dev")
	_ = es.Set("API_KEY", "secret")
	es.Delete("API_KEY")
	_, ok := es.Get("API_KEY")
	if ok {
		t.Error("expected key to be deleted")
	}
}

func TestMerge(t *testing.T) {
	base, _ := New("base")
	_ = base.Set("APP_ENV", "development")
	_ = base.Set("LOG_LEVEL", "debug")

	override, _ := New("override")
	_ = override.Set("LOG_LEVEL", "info")
	_ = override.Set("FEATURE_FLAG", "true")

	base.Merge(override)

	if v, _ := base.Get("LOG_LEVEL"); v != "info" {
		t.Errorf("expected LOG_LEVEL='info', got '%s'", v)
	}
	if _, ok := base.Get("FEATURE_FLAG"); !ok {
		t.Error("expected FEATURE_FLAG to exist after merge")
	}
	if v, _ := base.Get("APP_ENV"); v != "development" {
		t.Errorf("expected APP_ENV='development', got '%s'", v)
	}
}

func TestKeys(t *testing.T) {
	es, _ := New("dev")
	_ = es.Set("FOO", "1")
	_ = es.Set("BAR", "2")
	keys := es.Keys()
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
}
