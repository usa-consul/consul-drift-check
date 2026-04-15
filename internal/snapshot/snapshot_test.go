package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/consul-drift-check/internal/diff"
	"github.com/your-org/consul-drift-check/internal/snapshot"
)

func sampleResults() []diff.DiffResult {
	return []diff.DiffResult{
		{Key: "app/db", Status: diff.Modified, SourceValue: []byte("old"), DestValue: []byte("new")},
		{Key: "app/host", Status: diff.OnlyInSource},
	}
}

func TestSave_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	if err := snapshot.Save(path, "app/", sampleResults()); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("expected snapshot file to exist")
	}
}

func TestLoad_ReturnsSnapshot(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	results := sampleResults()
	if err := snapshot.Save(path, "app/", results); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	s, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if s.Prefix != "app/" {
		t.Errorf("expected prefix 'app/', got %q", s.Prefix)
	}

	if len(s.Results) != len(results) {
		t.Errorf("expected %d results, got %d", len(results), len(s.Results))
	}

	if s.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}

	if s.Timestamp.After(time.Now().Add(time.Second)) {
		t.Error("timestamp appears to be in the future")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")

	if err := os.WriteFile(path, []byte("not-json{"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := snapshot.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}
