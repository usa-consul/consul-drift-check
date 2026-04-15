package filter_test

import (
	"testing"

	"github.com/hashicorp/consul/api"

	"github.com/your-org/consul-drift-check/internal/filter"
)

func makeEntries(keys ...string) api.KVPairs {
	pairs := make(api.KVPairs, 0, len(keys))
	for _, k := range keys {
		pairs = append(pairs, &api.KVPair{Key: k, Value: []byte("v")})
	}
	return pairs
}

func keys(pairs api.KVPairs) []string {
	out := make([]string, 0, len(pairs))
	for _, p := range pairs {
		out = append(out, p.Key)
	}
	return out
}

func TestApply_NoPrefixFilter_ReturnsAll(t *testing.T) {
	entries := makeEntries("app/a", "app/b", "db/c")
	result := filter.Apply(entries, filter.Options{})
	if len(result) != 3 {
		t.Errorf("expected 3 entries, got %d", len(result))
	}
}

func TestApply_IncludePrefixes_FiltersCorrectly(t *testing.T) {
	entries := makeEntries("app/a", "app/b", "db/c", "infra/d")
	opts := filter.Options{Prefixes: []string{"app/", "infra/"}}
	result := filter.Apply(entries, opts)
	if len(result) != 3 {
		t.Errorf("expected 3 entries, got %d: %v", len(result), keys(result))
	}
}

func TestApply_ExcludeKeys_RemovesExactMatch(t *testing.T) {
	entries := makeEntries("app/secret", "app/config", "app/token")
	opts := filter.Options{ExcludeKeys: []string{"app/secret", "app/token"}}
	result := filter.Apply(entries, opts)
	if len(result) != 1 || result[0].Key != "app/config" {
		t.Errorf("expected only app/config, got %v", keys(result))
	}
}

func TestApply_ExcludePrefixes_RemovesMatchingKeys(t *testing.T) {
	entries := makeEntries("app/config", "app/secrets/key1", "app/secrets/key2", "db/host")
	opts := filter.Options{ExcludePrefixes: []string{"app/secrets/"}}
	result := filter.Apply(entries, opts)
	if len(result) != 2 {
		t.Errorf("expected 2 entries, got %d: %v", len(result), keys(result))
	}
}

func TestApply_CombinedFilters(t *testing.T) {
	entries := makeEntries("app/config", "app/secrets/key1", "db/host", "infra/vm")
	opts := filter.Options{
		Prefixes:        []string{"app/", "db/"},
		ExcludePrefixes: []string{"app/secrets/"},
	}
	result := filter.Apply(entries, opts)
	if len(result) != 2 {
		t.Errorf("expected 2 entries, got %d: %v", len(result), keys(result))
	}
}

func TestApply_NilEntrySkipped(t *testing.T) {
	entries := api.KVPairs{nil, {Key: "app/config", Value: []byte("v")}}
	result := filter.Apply(entries, filter.Options{})
	if len(result) != 1 {
		t.Errorf("expected 1 entry, got %d", len(result))
	}
}

func TestApply_EmptyInput_ReturnsEmpty(t *testing.T) {
	result := filter.Apply(api.KVPairs{}, filter.Options{})
	if len(result) != 0 {
		t.Errorf("expected 0 entries, got %d", len(result))
	}
}
