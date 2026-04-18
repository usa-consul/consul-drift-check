package prefix_test

import (
	"testing"

	"github.com/hashicorp/consul/api"

	"github.com/organisation/consul-drift-check/internal/prefix"
)

func makePairs(keys ...string) []*api.KVPair {
	pairs := make([]*api.KVPair, len(keys))
	for i, k := range keys {
		pairs[i] = &api.KVPair{Key: k, Value: []byte("v")}
	}
	return pairs
}

func TestAnalyse_EmptyInput(t *testing.T) {
	stats := prefix.Analyse(nil)
	if len(stats) != 0 {
		t.Fatalf("expected empty, got %d", len(stats))
	}
}

func TestAnalyse_GroupsByTopSegment(t *testing.T) {
	pairs := makePairs("app/db/host", "app/db/port", "infra/region", "app/name")
	stats := prefix.Analyse(pairs)
	if len(stats) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(stats))
	}
	if stats[0].Prefix != "app" {
		t.Errorf("expected first prefix 'app', got %q", stats[0].Prefix)
	}
	if stats[0].Count != 3 {
		t.Errorf("expected count 3, got %d", stats[0].Count)
	}
	if stats[1].Prefix != "infra" {
		t.Errorf("expected second prefix 'infra', got %q", stats[1].Prefix)
	}
}

func TestAnalyse_NilPairSkipped(t *testing.T) {
	pairs := []*api.KVPair{nil, {Key: "svc/a", Value: []byte("1")}}
	stats := prefix.Analyse(pairs)
	if len(stats) != 1 {
		t.Fatalf("expected 1 stat, got %d", len(stats))
	}
}

func TestAnalyse_LeadingSlashStripped(t *testing.T) {
	pairs := makePairs("/app/key")
	stats := prefix.Analyse(pairs)
	if stats[0].Prefix != "app" {
		t.Errorf("expected prefix 'app', got %q", stats[0].Prefix)
	}
}

func TestCommon_ReturnsSharedPrefixes(t *testing.T) {
	a := makePairs("app/x", "infra/y", "db/z")
	b := makePairs("app/a", "cache/b", "db/c")
	common := prefix.Common(a, b)
	if len(common) != 2 {
		t.Fatalf("expected 2 common prefixes, got %v", common)
	}
	if common[0] != "app" || common[1] != "db" {
		t.Errorf("unexpected common prefixes: %v", common)
	}
}

func TestCommon_NoOverlap_ReturnsEmpty(t *testing.T) {
	a := makePairs("alpha/x")
	b := makePairs("beta/y")
	if got := prefix.Common(a, b); len(got) != 0 {
		t.Errorf("expected empty, got %v", got)
	}
}
