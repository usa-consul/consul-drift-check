package score_test

import (
	"testing"

	"github.com/your-org/consul-drift-check/internal/diff"
	"github.com/your-org/consul-drift-check/internal/score"
)

func makeResults(statuses ...diff.Status) []diff.Result {
	out := make([]diff.Result, len(statuses))
	for i, s := range statuses {
		out[i] = diff.Result{Key: "k", Status: s}
	}
	return out
}

func TestCompute_EmptyResults(t *testing.T) {
	s := score.Compute(nil)
	if s.Score != 0 {
		t.Fatalf("expected 0, got %f", s.Score)
	}
}

func TestCompute_OnlyModified(t *testing.T) {
	s := score.Compute(makeResults(diff.StatusModified, diff.StatusModified))
	if s.Modified != 2 {
		t.Fatalf("expected 2 modified, got %d", s.Modified)
	}
	if s.Score != 4.0 {
		t.Fatalf("expected score 4.0, got %f", s.Score)
	}
}

func TestCompute_MixedStatuses(t *testing.T) {
	results := makeResults(
		diff.StatusModified,
		diff.StatusOnlyInSource,
		diff.StatusOnlyInDestination,
	)
	s := score.Compute(results)
	expected := score.WeightModified + score.WeightOnlyInSource + score.WeightOnlyInDest
	if s.Score != expected {
		t.Fatalf("expected %f, got %f", expected, s.Score)
	}
	if s.Modified != 1 || s.OnlyInSource != 1 || s.OnlyInDest != 1 {
		t.Fatalf("unexpected counts: %+v", s)
	}
}

func TestCompute_NoDrift(t *testing.T) {
	results := makeResults(diff.StatusEqual, diff.StatusEqual)
	s := score.Compute(results)
	if s.Score != 0 {
		t.Fatalf("expected 0 for equal results, got %f", s.Score)
	}
}
