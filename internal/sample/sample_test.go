package sample_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"

	"github.com/your-org/consul-drift-check/internal/sample"
)

func samplePairs() []*api.KVPair {
	return []*api.KVPair{
		{Key: "app/foo", Value: []byte("bar")},
		{Key: "app/baz", Value: []byte("qux")},
	}
}

func TestRecord_CreatesFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "samples.ndjson")
	if err := sample.Record(path, "app/", samplePairs()); err != nil {
		t.Fatalf("Record: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}

func TestLoadAll_ReturnsEntries(t *testing.T) {
	path := filepath.Join(t.TempDir(), "samples.ndjson")
	pairs := samplePairs()

	for i := 0; i < 3; i++ {
		if err := sample.Record(path, "app/", pairs); err != nil {
			t.Fatalf("Record: %v", err)
		}
	}

	entries, err := sample.LoadAll(path)
	if err != nil {
		t.Fatalf("LoadAll: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Prefix != "app/" {
		t.Errorf("unexpected prefix: %s", entries[0].Prefix)
	}
}

func TestLoadAll_MissingFile(t *testing.T) {
	entries, err := sample.LoadAll("/nonexistent/path/samples.ndjson")
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries, got %v", entries)
	}
}

func TestSince_FiltersEntries(t *testing.T) {
	now := time.Now().UTC()
	all := []sample.Entry{
		{Timestamp: now.Add(-2 * time.Hour), Prefix: "old"},
		{Timestamp: now.Add(-30 * time.Minute), Prefix: "recent"},
		{Timestamp: now, Prefix: "now"},
	}

	cutoff := now.Add(-1 * time.Hour)
	got := sample.Since(all, cutoff)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries after cutoff, got %d", len(got))
	}
	if got[0].Prefix != "recent" {
		t.Errorf("unexpected first entry: %s", got[0].Prefix)
	}
}
