package audit_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/example/consul-drift-check/internal/audit"
	"github.com/example/consul-drift-check/internal/metrics"
)

func makeSummary(added, removed, modified int) metrics.Summary {
	return metrics.Summary{
		Added:    added,
		Removed:  removed,
		Modified: modified,
		Total:    added + removed + modified,
	}
}

func TestRecord_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.jsonl")

	l := audit.NewLogger(path)
	err := l.Record(audit.Entry{
		ConfigPath: "config.yaml",
		Prefix:     "app/",
		Summary:    makeSummary(1, 0, 0),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("log file not created: %v", err)
	}
}

func TestRecord_AppendsMultipleEntries(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.jsonl")
	l := audit.NewLogger(path)

	for i := 0; i < 3; i++ {
		if err := l.Record(audit.Entry{ConfigPath: "cfg", Summary: makeSummary(i, 0, 0)}); err != nil {
			t.Fatalf("record %d: %v", i, err)
		}
	}

	entries, err := audit.ReadAll(path)
	if err != nil {
		t.Fatalf("read all: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestRecord_TimestampSetAutomatically(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.jsonl")
	l := audit.NewLogger(path)

	before := time.Now().UTC()
	if err := l.Record(audit.Entry{ConfigPath: "cfg"}); err != nil {
		t.Fatalf("record: %v", err)
	}

	entries, _ := audit.ReadAll(path)
	if entries[0].Timestamp.Before(before) {
		t.Error("timestamp should be set to current time")
	}
}

func TestReadAll_MissingFile(t *testing.T) {
	_, err := audit.ReadAll("/nonexistent/audit.jsonl")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
