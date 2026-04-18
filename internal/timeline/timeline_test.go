package timeline_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/consul-drift-check/internal/diff"
	"github.com/your-org/consul-drift-check/internal/timeline"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "app/env", Status: "modified", SourceValue: []byte("prod"), DestValue: []byte("staging")},
	}
}

func TestRecord_CreatesFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tl.json")
	if err := timeline.Record(path, sampleResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not created: %v", err)
	}
}

func TestLoad_ReturnsEntries(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tl.json")
	if err := timeline.Record(path, sampleResults()); err != nil {
		t.Fatal(err)
	}
	tl, err := timeline.Load(path)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if len(tl) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(tl))
	}
	if len(tl[0].Results) != 1 {
		t.Errorf("expected 1 result, got %d", len(tl[0].Results))
	}
}

func TestRecord_AppendsEntries(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tl.json")
	for i := 0; i < 3; i++ {
		if err := timeline.Record(path, sampleResults()); err != nil {
			t.Fatal(err)
		}
	}
	tl, err := timeline.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(tl) != 3 {
		t.Errorf("expected 3 entries, got %d", len(tl))
	}
}

func TestSince_FiltersEntries(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tl.json")
	if err := timeline.Record(path, sampleResults()); err != nil {
		t.Fatal(err)
	}
	tl, _ := timeline.Load(path)

	got := tl.Since(time.Now().Add(time.Hour))
	if len(got) != 0 {
		t.Errorf("expected 0 entries after future time, got %d", len(got))
	}

	got = tl.Since(time.Now().Add(-time.Hour))
	if len(got) != 1 {
		t.Errorf("expected 1 entry, got %d", len(got))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := timeline.Load("/nonexistent/tl.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
