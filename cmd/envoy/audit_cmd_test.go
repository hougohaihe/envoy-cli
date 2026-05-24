package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/envoy-cli/internal/audit"
)

func setupAuditLog(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.jsonl")

	logger, err := audit.NewLogger(path)
	if err != nil {
		t.Fatalf("NewLogger: %v", err)
	}

	e1 := audit.NewEvent(audit.EventPush, "staging", true)
	e1.Message = "pushed ok"
	e2 := audit.NewEvent(audit.EventPull, "production", false)
	e2.Message = "remote error"

	if err := logger.Log(e1); err != nil {
		t.Fatalf("Log: %v", err)
	}
	if err := logger.Log(e2); err != nil {
		t.Fatalf("Log: %v", err)
	}
	return path
}

func TestAuditListCmd_TableOutput(t *testing.T) {
	path := setupAuditLog(t)
	cmd := buildAuditCmd(path)

	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"list"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "push") {
		t.Errorf("expected 'push' in output, got: %s", out)
	}
	if !strings.Contains(out, "staging") {
		t.Errorf("expected 'staging' in output, got: %s", out)
	}
}

func TestAuditListCmd_JSONOutput(t *testing.T) {
	path := setupAuditLog(t)
	cmd := buildAuditCmd(path)

	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"list", "--json"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !strings.HasPrefix(strings.TrimSpace(buf.String()), "[") {
		t.Errorf("expected JSON array, got: %s", buf.String())
	}
}

func TestAuditListCmd_FilterEnvSet(t *testing.T) {
	path := setupAuditLog(t)
	cmd := buildAuditCmd(path)

	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"list", "--env-set", "staging"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "production") {
		t.Errorf("expected 'production' to be filtered out, got: %s", out)
	}
	if !strings.Contains(out, "staging") {
		t.Errorf("expected 'staging' in filtered output, got: %s", out)
	}
}
