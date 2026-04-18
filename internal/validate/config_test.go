package validate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholasgasior/consul-drift-check/internal/validate"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "validate-*.yaml")
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
	cfg, err := validate.LoadConfig("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Rules) != 0 {
		t.Errorf("expected empty rules")
	}
}

func TestLoadConfig_ValidFile(t *testing.T) {
	path := writeTemp(t, "rules:\n  - prefix: app/\n    required: true\n    severity: error\n")
	cfg, err := validate.LoadConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(cfg.Rules))
	}
	if cfg.Rules[0].Prefix != "app/" {
		t.Errorf("unexpected prefix: %s", cfg.Rules[0].Prefix)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.yaml")
	_, err := validate.LoadConfig(path)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	path := writeTemp(t, ": : invalid")
	_, err := validate.LoadConfig(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}
