package summarize_test

import (
	"strings"
	"testing"

	"github.com/your-org/consul-drift-check/internal/diff"
	"github.com/your-org/consul-drift-check/internal/summarize"
)

func makeResults(statuses ...diff.Status) []diff.Result {
	results := make([]diff.Result, len(statuses))
	for i, s := range statuses {
		results[i] = diff.Result{Key: fmt.Sprintf("key/%d", i), Status: s}
	}
	return results
}

func TestBuild_EmptyResults(t *testing.T) {
	s := summarize.Build(nil)
	if s.Total != 0 || s.Added != 0 || s.Removed != 0 || s.Modified != 0 {
		t.Fatalf("expected zero summary, got %+v", s)
	}
}

func TestBuild_CountsCorrectly(t *testing.T) {
	results := []diff.Result{
		{Key: "a", Status: diff.StatusOnlyInSource},
		{Key: "b", Status: diff.StatusOnlyInDestination},
		{Key: "c", Status: diff.StatusModified},
		{Key: "d", Status: diff.StatusModified},
	}
	s := summarize.Build(results)
	if s.Total != 4 { t.Errorf("total: want 4 got %d", s.Total) }
	if s.Added != 1 { t.Errorf("added: want 1 got %d", s.Added) }
	if s.Removed != 1 { t.Errorf("removed: want 1 got %d", s.Removed) }
	if s.Modified != 2 { t.Errorf("modified: want 2 got %d", s.Modified) }
}

func TestReport_NoDrift(t *testing.T) {
	s := summarize.Build(nil)
	if s.Report() != "no drift detected" {
		t.Errorf("unexpected report: %q", s.Report())
	}
}

func TestReport_ContainsPrefix(t *testing.T) {
	results := []diff.Result{
		{Key: "foo", Status: diff.StatusOnlyInSource},
		{Key: "bar", Status: diff.StatusOnlyInDestination},
		{Key: "baz", Status: diff.StatusModified},
	}
	s := summarize.Build(results)
	report := s.Report()
	for _, prefix := range []string{"+", "-", "~"} {
		if !strings.Contains(report, prefix) {
			t.Errorf("report missing prefix %q", prefix)
		}
	}
}

func TestString_Format(t *testing.T) {
	s := summarize.Summary{Total: 3, Added: 1, Removed: 1, Modified: 1}
	got := s.String()
	if !strings.Contains(got, "total=3") {
		t.Errorf("String() missing total: %q", got)
	}
}
