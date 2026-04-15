package metrics

import (
	"testing"
	"time"

	"github.com/your-org/consul-drift-check/internal/diff"
)

func makeResults(statuses []diff.DriftStatus) []diff.Result {
	results := make([]diff.Result, len(statuses))
	for i, s := range statuses {
		results[i] = diff.Result{Key: "key", Status: s}
	}
	return results
}

func TestCollect_NoDrift(t *testing.T) {
	results := makeResults([]diff.DriftStatus{
		diff.StatusMatch, diff.StatusMatch,
	})
	s := Collect(results, 10*time.Millisecond)

	if s.TotalKeys != 2 {
		t.Errorf("expected TotalKeys=2, got %d", s.TotalKeys)
	}
	if s.MatchedKeys != 2 {
		t.Errorf("expected MatchedKeys=2, got %d", s.MatchedKeys)
	}
	if s.DriftDetected {
		t.Error("expected DriftDetected=false")
	}
}

func TestCollect_WithDrift(t *testing.T) {
	results := makeResults([]diff.DriftStatus{
		diff.StatusMatch,
		diff.StatusOnlyInSource,
		diff.StatusOnlyInDest,
		diff.StatusModified,
	})
	s := Collect(results, 50*time.Millisecond)

	if s.TotalKeys != 4 {
		t.Errorf("expected TotalKeys=4, got %d", s.TotalKeys)
	}
	if s.OnlyInSource != 1 {
		t.Errorf("expected OnlyInSource=1, got %d", s.OnlyInSource)
	}
	if s.OnlyInDest != 1 {
		t.Errorf("expected OnlyInDest=1, got %d", s.OnlyInDest)
	}
	if s.ModifiedKeys != 1 {
		t.Errorf("expected ModifiedKeys=1, got %d", s.ModifiedKeys)
	}
	if !s.DriftDetected {
		t.Error("expected DriftDetected=true")
	}
}

func TestCollect_EmptyResults(t *testing.T) {
	s := Collect([]diff.Result{}, 0)

	if s.TotalKeys != 0 {
		t.Errorf("expected TotalKeys=0, got %d", s.TotalKeys)
	}
	if s.DriftDetected {
		t.Error("expected DriftDetected=false for empty results")
	}
	if s.CheckedAt.IsZero() {
		t.Error("expected CheckedAt to be set")
	}
}
