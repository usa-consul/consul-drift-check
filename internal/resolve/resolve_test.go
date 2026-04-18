package resolve_test

import (
	"testing"

	"github.com/hashicorp/consul/api"

	"github.com/organisation/consul-drift-check/internal/resolve"
)

func makePairs(keys ...string) []*api.KVPair {
	pairs := make([]*api.KVPair, len(keys))
	for i, k := range keys {
		pairs[i] = &api.KVPair{Key: k, Value: []byte("v")}
	}
	return pairs
}

func TestApply_NoRules_FallbackToKey(t *testing.T) {
	pairs := makePairs("service/web/port")
	opts := resolve.Options{FallbackToKey: true}
	results := resolve.Apply(pairs, opts)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Alias != "service/web/port" {
		t.Errorf("expected alias %q, got %q", "service/web/port", results[0].Alias)
	}
}

func TestApply_NoRules_NoFallback_EmptyAlias(t *testing.T) {
	pairs := makePairs("service/web/port")
	opts := resolve.Options{FallbackToKey: false}
	results := resolve.Apply(pairs, opts)
	if results[0].Alias != "" {
		t.Errorf("expected empty alias, got %q", results[0].Alias)
	}
}

func TestApply_MatchingPrefix_AssignsAlias(t *testing.T) {
	pairs := makePairs("service/web/port", "service/db/host")
	opts := resolve.Options{
		Rules: []resolve.Rule{
			{Prefix: "service/web", Alias: "web-service"},
		},
		FallbackToKey: true,
	}
	results := resolve.Apply(pairs, opts)
	if results[0].Alias != "web-service/port" {
		t.Errorf("expected %q, got %q", "web-service/port", results[0].Alias)
	}
	if results[1].Alias != "service/db/host" {
		t.Errorf("expected fallback %q, got %q", "service/db/host", results[1].Alias)
	}
}

func TestApply_ExactKeyMatch_ReturnsAlias(t *testing.T) {
	pairs := makePairs("config/timeout")
	opts := resolve.Options{
		Rules: []resolve.Rule{
			{Prefix: "config/timeout", Alias: "timeout-setting"},
		},
	}
	results := resolve.Apply(pairs, opts)
	if results[0].Alias != "timeout-setting" {
		t.Errorf("expected %q, got %q", "timeout-setting", results[0].Alias)
	}
}

func TestApply_NilPair_Skipped(t *testing.T) {
	pairs := []*api.KVPair{nil, {Key: "k", Value: []byte("v")}}
	opts := resolve.Options{FallbackToKey: true}
	results := resolve.Apply(pairs, opts)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestApply_FirstRuleWins(t *testing.T) {
	pairs := makePairs("svc/api/v1")
	opts := resolve.Options{
		Rules: []resolve.Rule{
			{Prefix: "svc/api", Alias: "api"},
			{Prefix: "svc", Alias: "generic"},
		},
	}
	results := resolve.Apply(pairs, opts)
	if results[0].Alias != "api/v1" {
		t.Errorf("expected %q, got %q", "api/v1", results[0].Alias)
	}
}
