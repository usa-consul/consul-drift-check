package groupby_test

import (
	"testing"

	"github.com/your-org/consul-drift-check/internal/diff"
	"github.com/your-org/consul-drift-check/internal/groupby"
)

func makeResults(keys ...string) []diff.Result {
	out := make([]diff.Result, len(keys))
	for i, k := range keys {
		out[i] = diff.Result{Key: k, Status: diff.Modified}
	}
	return out
}

func TestApply_GroupsByFirstSegment(t *testing.T) {
	results := makeResults("services/web/port", "services/api/port", "infra/db/host")
	g := groupby.Apply(results, groupby.Options{SegmentIndex: 0})

	if len(g["services"]) != 2 {
		t.Fatalf("expected 2 in services, got %d", len(g["services"]))
	}
	if len(g["infra"]) != 1 {
		t.Fatalf("expected 1 in infra, got %d", len(g["infra"]))
	}
}

func TestApply_GroupsBySecondSegment(t *testing.T) {
	results := makeResults("services/web/port", "services/web/host", "services/api/port")
	g := groupby.Apply(results, groupby.Options{SegmentIndex: 1})

	if len(g["web"]) != 2 {
		t.Fatalf("expected 2 in web, got %d", len(g["web"]))
	}
	if len(g["api"]) != 1 {
		t.Fatalf("expected 1 in api, got %d", len(g["api"]))
	}
}

func TestApply_EmptyResults_ReturnsEmptyGroup(t *testing.T) {
	g := groupby.Apply(nil, groupby.Options{})
	if len(g) != 0 {
		t.Fatalf("expected empty group, got %d keys", len(g))
	}
}

func TestApply_IndexOutOfRange_FallsBackToOther(t *testing.T) {
	results := makeResults("shallow")
	g := groupby.Apply(results, groupby.Options{SegmentIndex: 5})
	if len(g["_other"]) != 1 {
		t.Fatalf("expected 1 in _other, got %d", len(g["_other"]))
	}
}

func TestKeys_ReturnsSorted(t *testing.T) {
	results := makeResults("z/key", "a/key", "m/key")
	g := groupby.Apply(results, groupby.Options{SegmentIndex: 0})
	keys := groupby.Keys(g)
	if keys[0] != "a" || keys[1] != "m" || keys[2] != "z" {
		t.Fatalf("unexpected order: %v", keys)
	}
}

func TestApply_CustomSeparator(t *testing.T) {
	results := makeResults("services.web.port", "services.api.port")
	g := groupby.Apply(results, groupby.Options{SegmentIndex: 0, Separator: "."})
	if len(g["services"]) != 2 {
		t.Fatalf("expected 2 in services, got %d", len(g["services"]))
	}
}
