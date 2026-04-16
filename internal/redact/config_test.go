package redact_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/consul-drift-check/internal/redact"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "redact-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestLoadConfig_EmptyPath_ReturnsEmpty(t *testing.T) {
	opts, err := redact.LoadConfig("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(opts.SensitiveKeys) != 0 {
		t.Fatal("expected no sensitive keys")
	}
}

func TestLoadConfig_ValidFile(t *testing.T) {
	path := writeTemp(t, "sensitive_keys:\n  - password\n  - token\nmask: '[HIDDEN]'\n")
	opts, err := redact.LoadConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(opts.SensitiveKeys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(opts.SensitiveKeys))
	}
	if opts.Mask != "[HIDDEN]" {
		t.Fatalf("expected mask '[HIDDEN]', got %s", opts.Mask)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := redact.LoadConfig(filepath.Join(t.TempDir(), "missing.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	path := writeTemp(t, ": invalid: [yaml")
	_, err := redact.LoadConfig(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}
