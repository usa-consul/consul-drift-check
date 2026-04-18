package diff

import (
	"testing"
	"time"
)

func makeTrendEntries() []HistoryEntry {
	now := time.Now().UTC()
	return []HistoryEntry{
		{
			Timestamp: now.Add(-2 * time.Hour),
			Results: []Result{
				{Key: "a", Status: "only_source"},
				{Key: "b", Status: "modified"},
			},
		},
		{
			Timestamp: now.Add(-1 * time.Hour),
			Results: []Result{
				{Key: "a", Status: "only_source"},
				{Key: "b", Status: "modified"},
				{Key: "c", Status: "only_destination"},
			},
		},
		{
			Timestamp: now,
			Results:   []Result{},
		},
	}
}

func TestBuildTrend_PointCount(t *testing.T) {
	entries := makeTrendEntries()
	trend := BuildTrend(entries)
	if len(trend.Points) != 3 {
		t.Fatalf("expected 3 points, got %d", len(trend.Points))
	}
}

func TestBuildTrend_CountsCorrectly(t *testing.T) {
	entries := makeTrendEntries()
	trend := BuildTrend(entries)
	p := trend.Points[1]
	if p.OnlySource != 1 || p.OnlyDest != 1 || p.Modified != 1 || p.Total != 3 {
		t.Errorf("unexpected counts: %+v", p)
	}
}

func TestBuildTrend_SortedAscending(t *testing.T) {
	entries := makeTrendEntries()
	// reverse order
	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
	}
	trend := BuildTrend(entries)
	for i := 1; i < len(trend.Points); i++ {
		if trend.Points[i].Timestamp.Before(trend.Points[i-1].Timestamp) {
			t.Errorf("points not sorted at index %d", i)
		}
	}
}

func TestDelta_PositiveGrowth(t *testing.T) {
	entries := makeTrendEntries()
	trend := BuildTrend(entries)
	delta := trend.Delta()
	if delta != -2 {
		t.Errorf("expected delta -2, got %d", delta)
	}
}

func TestDelta_EmptyTrend(t *testing.T) {
	trend := Trend{}
	if trend.Delta() != 0 {
		t.Error("expected 0 delta for empty trend")
	}
}

func TestBuildTrend_EmptyEntries(t *testing.T) {
	trend := BuildTrend(nil)
	if len(trend.Points) != 0 {
		t.Errorf("expected 0 points, got %d", len(trend.Points))
	}
}
