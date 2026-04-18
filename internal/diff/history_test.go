package diff

import (
	"path/filepath"
	"testing"
	"time"
)

func sampleEntry(prefix string, status string) HistoryEntry {
	return HistoryEntry{
		Timestamp: time.Now().UTC(),
		Prefix:    prefix,
		Results: []Result{
			{Key: "app/key", Status: status},
		},
	}
}

func TestAppendHistory_CreatesFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.ndjson")
	entry := sampleEntry("app/", StatusModified)
	if err := AppendHistory(path, entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries, err := LoadHistory(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("want 1 entry, got %d", len(entries))
	}
	if entries[0].Prefix != "app/" {
		t.Errorf("want prefix app/, got %s", entries[0].Prefix)
	}
}

func TestAppendHistory_AppendsMultiple(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.ndjson")
	for i := 0; i < 3; i++ {
		if err := AppendHistory(path, sampleEntry("svc/", StatusOnlyInSource)); err != nil {
			t.Fatalf("append %d: %v", i, err)
		}
	}
	entries, err := LoadHistory(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("want 3, got %d", len(entries))
	}
}

func TestLoadHistory_MissingFile(t *testing.T) {
	entries, err := LoadHistory("/nonexistent/path/history.ndjson")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil for missing file")
	}
}

func TestSince_FiltersOldEntries(t *testing.T) {
	now := time.Now().UTC()
	old := HistoryEntry{Timestamp: now.Add(-2 * time.Hour), Prefix: "old/"}
	recent := HistoryEntry{Timestamp: now.Add(-30 * time.Minute), Prefix: "recent/"}
	cutoff := now.Add(-1 * time.Hour)
	out := Since([]HistoryEntry{old, recent}, cutoff)
	if len(out) != 1 {
		t.Fatalf("want 1, got %d", len(out))
	}
	if out[0].Prefix != "recent/" {
		t.Errorf("unexpected prefix: %s", out[0].Prefix)
	}
}

func TestSince_EmptyInput(t *testing.T) {
	out := Since(nil, time.Now())
	if len(out) != 0 {
		t.Errorf("expected empty, got %d", len(out))
	}
}
