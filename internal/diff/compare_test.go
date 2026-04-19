package diff

import (
	"testing"
)

func makeCompareResults(tuples [][2]string) []Result {
	out := make([]Result, len(tuples))
	for i, t := range tuples {
		out[i] = Result{Key: t[0], Status: t[1]}
	}
	return out
}

func TestCompareByPrefix_EmptyResults(t *testing.T) {
	got := CompareByPrefix(nil, CompareOptions{})
	if len(got) != 0 {
		t.Fatalf("expected empty, got %d", len(got))
	}
}

func TestCompareByPrefix_GroupsByTopSegment(t *testing.T) {
	results := makeCompareResults([][2]string{
		{"app/db/host", "modified"},
		{"app/db/port", "only_in_source"},
		{"infra/net", "only_in_destination"},
	})
	got := CompareByPrefix(results, CompareOptions{CaseSensitive: true})
	if len(got) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(got))
	}
	if got[0].Prefix != "app" {
		t.Errorf("expected app, got %s", got[0].Prefix)
	}
	if got[0].Modified != 1 || got[0].Added != 1 {
		t.Errorf("unexpected app counts: %+v", got[0])
	}
	if got[1].Prefix != "infra" || got[1].Removed != 1 {
		t.Errorf("unexpected infra counts: %+v", got[1])
	}
}

func TestCompareByPrefix_CaseInsensitive(t *testing.T) {
	results := makeCompareResults([][2]string{
		{"APP/key", "modified"},
		{"app/other", "modified"},
	})
	got := CompareByPrefix(results, CompareOptions{CaseSensitive: false})
	if len(got) != 1 {
		t.Fatalf("expected 1 group, got %d", len(got))
	}
	if got[0].Modified != 2 {
		t.Errorf("expected 2 modified, got %d", got[0].Modified)
	}
}

func TestCompareByPrefix_IgnorePrefixes(t *testing.T) {
	results := makeCompareResults([][2]string{
		{"secret/key", "modified"},
		{"app/key", "only_in_source"},
	})
	got := CompareByPrefix(results, CompareOptions{
		CaseSensitive:  true,
		IgnorePrefixes: []string{"secret"},
	})
	if len(got) != 1 {
		t.Fatalf("expected 1 group, got %d", len(got))
	}
	if got[0].Prefix != "app" {
		t.Errorf("expected app, got %s", got[0].Prefix)
	}
}

func TestCompareByPrefix_SortedAlphabetically(t *testing.T) {
	results := makeCompareResults([][2]string{
		{"zebra/k", "modified"},
		{"alpha/k", "modified"},
		{"mango/k", "modified"},
	})
	got := CompareByPrefix(results, CompareOptions{CaseSensitive: true})
	order := []string{"alpha", "mango", "zebra"}
	for i, p := range order {
		if got[i].Prefix != p {
			t.Errorf("pos %d: expected %s got %s", i, p, got[i].Prefix)
		}
	}
}
