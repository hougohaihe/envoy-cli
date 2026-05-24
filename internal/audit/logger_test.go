package audit_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envoy-cli/internal/audit"
)

func tempLogPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "audit", "events.jsonl")
}

func TestLogger_LogAndReadAll(t *testing.T) {
	path := tempLogPath(t)
	l, err := audit.NewLogger(path)
	if err != nil {
		t.Fatalf("NewLogger: %v", err)
	}

	e1 := audit.NewEvent(audit.EventPush, "staging", true)
	e2 := audit.NewEvent(audit.EventPull, "production", false)

	if err := l.Log(e1); err != nil {
		t.Fatalf("Log e1: %v", err)
	}
	if err := l.Log(e2); err != nil {
		t.Fatalf("Log e2: %v", err)
	}

	events, err := l.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].Type != audit.EventPush || events[0].EnvSet != "staging" {
		t.Errorf("unexpected event[0]: %+v", events[0])
	}
	if events[1].Type != audit.EventPull || events[1].Success {
		t.Errorf("unexpected event[1]: %+v", events[1])
	}
}

func TestLogger_ReadAll_EmptyFile(t *testing.T) {
	path := tempLogPath(t)
	l, err := audit.NewLogger(path)
	if err != nil {
		t.Fatalf("NewLogger: %v", err)
	}
	events, err := l.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll on missing file: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("expected 0 events, got %d", len(events))
	}
}

func TestLogger_CreatesParentDirs(t *testing.T) {
	path := tempLogPath(t)
	_, err := audit.NewLogger(path)
	if err != nil {
		t.Fatalf("NewLogger: %v", err)
	}
	if _, err := os.Stat(filepath.Dir(path)); err != nil {
		t.Errorf("parent dir not created: %v", err)
	}
}

func TestNewEvent_Fields(t *testing.T) {
	e := audit.NewEvent(audit.EventSet, "dev", true)
	if e.Type != audit.EventSet {
		t.Errorf("expected type %q, got %q", audit.EventSet, e.Type)
	}
	if e.EnvSet != "dev" {
		t.Errorf("expected envset %q, got %q", "dev", e.EnvSet)
	}
	if !e.Success {
		t.Error("expected success=true")
	}
	if e.ID == "" {
		t.Error("expected non-empty ID")
	}
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}
