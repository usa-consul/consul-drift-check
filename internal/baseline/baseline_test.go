package baseline_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/consul-drift-check/internal/baseline"
)

func sampleEntries() map[string][]byte {
	return map[string][]byte{
		"app/config/host": []byte("localhost"),
		"app/config/port": []byte("8080"),
	}
}

func TestSave_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	err := baseline.Save(path, "app/config", sampleEntries())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("expected file to exist")
	}
}

func TestLoad_ReturnsBaseline(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	entries := sampleEntries()
	if err := baseline.Save(path, "app/config", entries); err != nil {
		t.Fatalf("save: %v", err)
	}

	b, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if b.Prefix != "app/config" {
		t.Errorf("expected prefix app/config, got %s", b.Prefix)
	}
	if b.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
	if b.CreatedAt.After(time.Now().Add(time.Second)) {
		t.Error("CreatedAt is in the future")
	}
	if string(b.Entries["app/config/host"]) != "localhost" {
		t.Errorf("unexpected entry value: %s", b.Entries["app/config/host"])
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := baseline.Load("/nonexistent/baseline.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("not-json"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	_, err := baseline.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestCompare_DetectsDrift(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	original := sampleEntries()
	if err := baseline.Save(path, "app/config", original); err != nil {
		t.Fatalf("save: %v", err)
	}

	b, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	current := map[string][]byte{
		"app/config/host": []byte("remotehost"),
		"app/config/port": []byte("8080"),
	}

	results := baseline.Compare(b, current)
	if len(results) != 1 {
		t.Fatalf("expected 1 drift result, got %d", len(results))
	}
	if results[0].Key != "app/config/host" {
		t.Errorf("expected drift on app/config/host, got %s", results[0].Key)
	}
}
