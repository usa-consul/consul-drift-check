package diff

import (
	"testing"
)

func makeScopeResults(keys []string, statuses []string) []Result {
	out := make([]Result, len(keys))
	for i, k := range keys {
		out[i] = Result{Key: k, Status: statuses[i%len(statuses)]}
	}
	return out
}

func TestScope_EmptyInput_ReturnsNil(t *testing.T) {
	got := Scope(nil, ScopeOptions{})
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestScope_NoOptions_ReturnsAllSorted(t *testing.T) {
	results := makeScopeResults([]string{"z/key", "a/key", "m/key"}, []string{"modified"})
	got := Scope(results, ScopeOptions{})
	if len(got) != 3 {
		t.Fatalf("expected 3 results, got %d", len(got))
	}
	if got[0].Key != "a/key" || got[1].Key != "m/key" || got[2].Key != "z/key" {
		t.Fatalf("unexpected order: %v", got)
	}
}

func TestScope_PrefixFilter_LimitsResults(t *testing.T) {
	results := makeScopeResults(
		[]string{"service/a", "service/b", "config/c"},
		[]string{"modified"},
	)
	got := Scope(results, ScopeOptions{Prefixes: []string{"service/"}})
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
	for _, r := range got {
		if r.Key[:8] != "service/" {
			t.Fatalf("unexpected key outside prefix: %s", r.Key)
		}
	}
}

func TestScope_StatusFilter_LimitsResults(t *testing.T) {
	results := []Result{
		{Key: "a", Status: "modified"},
		{Key: "b", Status: "only-in-source"},
		{Key: "c", Status: "modified"},
	}
	got := Scope(results, ScopeOptions{Statuses: []string{"modified"}})
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestScope_StatusFilter_CaseInsensitive(t *testing.T) {
	results := []Result{
		{Key: "a", Status: "Modified"},
		{Key: "b", Status: "MODIFIED"},
		{Key: "c", Status: "only-in-source"},
	}
	got := Scope(results, ScopeOptions{Statuses: []string{"modified"}})
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestScope_MaxResults_CapsOutput(t *testing.T) {
	results := makeScopeResults(
		[]string{"a", "b", "c", "d", "e"},
		[]string{"modified"},
	)
	got := Scope(results, ScopeOptions{MaxResults: 3})
	if len(got) != 3 {
		t.Fatalf("expected 3, got %d", len(got))
	}
}

func TestScope_CombinedOptions(t *testing.T) {
	results := []Result{
		{Key: "svc/a", Status: "modified"},
		{Key: "svc/b", Status: "only-in-source"},
		{Key: "svc/c", Status: "modified"},
		{Key: "cfg/d", Status: "modified"},
	}
	got := Scope(results, ScopeOptions{
		Prefixes:   []string{"svc/"},
		Statuses:   []string{"modified"},
		MaxResults: 1,
	})
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
	if got[0].Key != "svc/a" {
		t.Fatalf("unexpected key: %s", got[0].Key)
	}
}
