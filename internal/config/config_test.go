package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/consul-drift-check/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}
	return p
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTemp(t, `
source:
  address: "http://consul-dc1:8500"
  token: "token-a"
  dc: "dc1"
destination:
  address: "http://consul-dc2:8500"
  token: "token-b"
  dc: "dc2"
prefix: "app/config"
format: "json"
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Source.Address != "http://consul-dc1:8500" {
		t.Errorf("expected source address, got %q", cfg.Source.Address)
	}
	if cfg.Format != "json" {
		t.Errorf("expected format json, got %q", cfg.Format)
	}
}

func TestLoad_EmptyPath(t *testing.T) {
	_, err := config.Load("")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestLoad_MissingSourceAddress(t *testing.T) {
	path := writeTemp(t, `
source:
  address: ""
destination:
  address: "http://consul-dc2:8500"
prefix: "app/config"
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing source address")
	}
}

func TestLoad_InvalidFormat(t *testing.T) {
	path := writeTemp(t, `
source:
  address: "http://consul-dc1:8500"
destination:
  address: "http://consul-dc2:8500"
prefix: "app/config"
format: "xml"
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for invalid format")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
