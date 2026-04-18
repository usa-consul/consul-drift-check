package overlap_test

import (
	"testing"

	"github.com/hashicorp/consul/api"

	"github.com/organisation/consul-drift-check/internal/overlap"
)

func makePairs(keys ...string) []*api.KVPair {
	out := make([]*api.KVPair, len(keys))
	for i, k := range keys {
		out[i] = &api.KVPair{Key: k, Value: []byte("v")}
	}
	return out
}

func TestFind_EmptyInput_ReturnsNil(t *testing.T) {
	if got := overlap.Find(nil, []string{"a/", "b/"}); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestFind_NoPrefixes_ReturnsNil(t *testing.T) {
	pairs := makePairs("a/key")
	if got := overlap.Find(pairs, nil); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestFind_NoOverlap_ReturnsEmpty(t *testing.T) {
	pairs := makePairs("prod/svc", "staging/svc")
	results := overlap.Find(pairs, []string{"prod/", "staging/"})
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestFind_OverlappingKey_Detected(t *testing.T) {
	pairs := makePairs("shared/config", "prod/only")
	prefixes := []string{"shared/", "prod/", "shared/"}
	// duplicate prefix intentional — should dedupe matches
	results := overlap.Find(pairs, prefixes)
	// shared/config matches "shared/" twice but should only appear once per prefix
	if len(results) != 0 {
		// same prefix twice should not count as overlap
		t.Fatalf("duplicate prefix should not produce overlap, got %v", results)
	}
}

func TestFind_KeyUnderTwoDifferentPrefixes_IsReturned(t *testing.T) {
	// craft a key that genuinely matches two distinct prefixes
	pairs := makePairs("a/b/key", "a/key", "c/key")
	results := overlap.Find(pairs, []string{"a/b/", "a/"})
	if len(results) != 1 {
		t.Fatalf("expected 1 overlap result, got %d", len(results))
	}
	if results[0].Key != "a/b/key" {
		t.Errorf("unexpected key %s", results[0].Key)
	}
	if len(results[0].Prefixes) != 2 {
		t.Errorf("expected 2 prefixes, got %v", results[0].Prefixes)
	}
}

func TestKeys_ExtractsKeyStrings(t *testing.T) {
	results := []overlap.Result{
		{Key: "z/key", Prefixes: []string{"z/"}},
		{Key: "a/key", Prefixes: []string{"a/"}},
	}
	keys := overlap.Keys(results)
	if len(keys) != 2 || keys[0] != "z/key" || keys[1] != "a/key" {
		t.Errorf("unexpected keys %v", keys)
	}
}
