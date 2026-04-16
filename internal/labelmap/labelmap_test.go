package labelmap_test

import (
	"testing"

	"github.com/hashicorp/consul/api"
	"github.com/nicholasgasior/consul-drift-check/internal/labelmap"
)

func makePairs(keys ...string) []*api.KVPair {
	pairs := make([]*api.KVPair, len(keys))
	for i, k := range keys {
		pairs[i] = &api.KVPair{Key: k, Value: []byte("v")}
	}
	return pairs
}

func TestApply_NoRules_ReturnsEmptyLabels(t *testing.T) {
	results := labelmap.Apply(makePairs("app/config"), nil)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if len(results[0].Labels) != 0 {
		t.Errorf("expected empty labels, got %v", results[0].Labels)
	}
}

func TestApply_MatchingPrefix_AssignsLabels(t *testing.T) {
	rules := []labelmap.Rule{
		{Prefix: "app/", Labels: map[string]string{"service": "app", "env": "prod"}},
	}
	results := labelmap.Apply(makePairs("app/config", "app/db"), rules)
	for _, r := range results {
		if r.Labels["service"] != "app" {
			t.Errorf("key %s: expected service=app, got %q", r.Key, r.Labels["service"])
		}
	}
}

func TestApply_FirstRuleWins(t *testing.T) {
	rules := []labelmap.Rule{
		{Prefix: "app/db", Labels: map[string]string{"service": "db"}},
		{Prefix: "app/", Labels: map[string]string{"service": "app"}},
	}
	results := labelmap.Apply(makePairs("app/db/host"), rules)
	if results[0].Labels["service"] != "db" {
		t.Errorf("expected first rule to win, got %q", results[0].Labels["service"])
	}
}

func TestApply_NoMatchingPrefix_EmptyLabels(t *testing.T) {
	rules := []labelmap.Rule{
		{Prefix: "infra/", Labels: map[string]string{"team": "ops"}},
	}
	results := labelmap.Apply(makePairs("app/config"), rules)
	if len(results[0].Labels) != 0 {
		t.Errorf("expected no labels, got %v", results[0].Labels)
	}
}

func TestApply_LabelMapIsCopied(t *testing.T) {
	rules := []labelmap.Rule{
		{Prefix: "app/", Labels: map[string]string{"env": "prod"}},
	}
	results := labelmap.Apply(makePairs("app/x", "app/y"), rules)
	results[0].Labels["env"] = "mutated"
	if results[1].Labels["env"] != "prod" {
		t.Errorf("label maps should be independent copies")
	}
}
