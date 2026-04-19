package diff

import (
	"testing"
)

func makeClusterResults(entries []struct {
	key    string
	status string
}) []Result {
	out := make([]Result, len(entries))
	for i, e := range entries {
		out[i] = Result{Key: e.key, Status: e.status}
	}
	return out
}

func TestBuildCluster_EmptyResults(t *testing.T) {
	s := BuildCluster("dc1", "dc2", nil)
	if len(s.Pairs) != 0 {
		t.Fatalf("expected 0 pairs, got %d", len(s.Pairs))
	}
}

func TestBuildCluster_PairLabel(t *testing.T) {
	results := makeClusterResults([]struct {
		key    string
		status string
	}{
		{"app/db", StatusModified},
	})
	s := BuildCluster("dc1", "dc2", results)
	if len(s.Pairs) != 1 {
		t.Fatalf("expected 1 pair, got %d", len(s.Pairs))
	}
	if s.Pairs[0].Pair != "dc1->dc2" {
		t.Errorf("unexpected pair label: %s", s.Pairs[0].Pair)
	}
}

func TestBuildCluster_CountsCorrectly(t *testing.T) {
	results := makeClusterResults([]struct {
		key    string
		status string
	}{
		{"app/key1", StatusOnlyInSource},
		{"app/key2", StatusOnlyInDestination},
		{"app/key3", StatusModified},
		{"app/key4", StatusModified},
	})
	s := BuildCluster("dc1", "dc2", results)
	if len(s.Pairs) != 1 {
		t.Fatalf("expected 1 segment group, got %d", len(s.Pairs))
	}
	p := s.Pairs[0]
	if p.Added != 1 {
		t.Errorf("Added: want 1, got %d", p.Added)
	}
	if p.Removed != 1 {
		t.Errorf("Removed: want 1, got %d", p.Removed)
	}
	if p.Modified != 2 {
		t.Errorf("Modified: want 2, got %d", p.Modified)
	}
	if p.Total != 4 {
		t.Errorf("Total: want 4, got %d", p.Total)
	}
}

func TestBuildCluster_MultipleSegments_SortedByPair(t *testing.T) {
	results := makeClusterResults([]struct {
		key    string
		status string
	}{
		{"svc-b/key", StatusModified},
		{"svc-a/key", StatusOnlyInSource},
	})
	s := BuildCluster("dc1", "dc2", results)
	if len(s.Pairs) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(s.Pairs))
	}
	if s.Pairs[0].Pair > s.Pairs[1].Pair {
		t.Error("pairs not sorted ascending")
	}
}

func TestBuildCluster_LeadingSlashStripped(t *testing.T) {
	results := makeClusterResults([]struct {
		key    string
		status string
	}{
		{"/app/key", StatusModified},
	})
	s := BuildCluster("dc1", "dc2", results)
	if len(s.Pairs) != 1 {
		t.Fatalf("expected 1 pair, got %d", len(s.Pairs))
	}
}
