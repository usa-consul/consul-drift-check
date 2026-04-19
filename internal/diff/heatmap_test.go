package diff

import (
	"testing"
	"time"
)

func makeHeatEntries() []HistoryEntry {
	now := time.Now()
	return []HistoryEntry{
		{
			Timestamp: now.Add(-2 * time.Hour),
			Results: []Result{
				{Key: "app/db/host", Status: "modified"},
				{Key: "app/db/port", Status: "modified"},
				{Key: "infra/vpc", Status: "only_in_source"},
			},
		},
		{
			Timestamp: now.Add(-1 * time.Hour),
			Results: []Result{
				{Key: "app/cache/ttl", Status: "modified"},
				{Key: "infra/subnet", Status: "only_in_destination"},
			},
		},
	}
}

func TestBuildHeatMap_EmptyEntries(t *testing.T) {
	hm := BuildHeatMap(nil)
	if len(hm) != 0 {
		t.Fatalf("expected empty heatmap, got %d cells", len(hm))
	}
}

func TestBuildHeatMap_GroupsByTopSegment(t *testing.T) {
	hm := BuildHeatMap(makeHeatEntries())
	if len(hm) != 2 {
		t.Fatalf("expected 2 cells, got %d", len(hm))
	}
	if hm[0].Prefix != "app" {
		t.Errorf("expected first prefix to be 'app', got %q", hm[0].Prefix)
	}
}

func TestBuildHeatMap_ScoresCorrectly(t *testing.T) {
	hm := BuildHeatMap(makeHeatEntries())
	var appCell HeatCell
	for _, c := range hm {
		if c.Prefix == "app" {
			appCell = c
		}
	}
	// 3 modified keys * 2.0 = 6.0
	if appCell.Score != 6.0 {
		t.Errorf("expected app score 6.0, got %f", appCell.Score)
	}
	if appCell.Count != 3 {
		t.Errorf("expected app count 3, got %d", appCell.Count)
	}
}

func TestBuildHeatMap_SortedDescending(t *testing.T) {
	hm := BuildHeatMap(makeHeatEntries())
	for i := 1; i < len(hm); i++ {
		if hm[i].Score > hm[i-1].Score {
			t.Errorf("heatmap not sorted descending at index %d", i)
		}
	}
}

func TestBuildHeatMap_LeadingSlashStripped(t *testing.T) {
	entries := []HistoryEntry{
		{Results: []Result{{Key: "/service/config", Status: "modified"}}},
	}
	hm := BuildHeatMap(entries)
	if len(hm) == 0 {
		t.Fatal("expected at least one cell")
	}
	if hm[0].Prefix != "service" {
		t.Errorf("expected prefix 'service', got %q", hm[0].Prefix)
	}
}
