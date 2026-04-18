package annotate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/example/consul-drift-check/internal/annotate"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "annotate-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestLoadConfig_EmptyPath_ReturnsNil(t *testing.T) {
	rules, err := annotate.LoadConfig("")
	if err != nil || rules != nil {
		t.Errorf("expected nil rules and no error, got %v %v", rules, err)
	}
}

func TestLoadConfig_ValidFile(t *testing.T) {
	path := writeTemp(t, "rules:\n  - prefix: app/\n    annotations:\n      team: platform\n")
	rules, err := annotate.LoadConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 || rules[0].Prefix != "app/" {
		t.Errorf("unexpected rules: %+v", rules)
	}
	if rules[0].Annotations["team"] != "platform" {
		t.Errorf("unexpected annotation value")
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := annotate.LoadConfig(filepath.Join(t.TempDir(), "missing.yaml"))
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	path := writeTemp(t, ": : invalid")
	_, err := annotate.LoadConfig(path)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}
