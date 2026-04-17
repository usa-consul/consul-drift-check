package dedupe_test

import (
	"testing"

	"github.com/hashicorp/consul/api"

	"github.com/example/consul-drift-check/internal/dedupe"
)

func makePairs(kvs ...string) []*api.KVPair {
	if len(kvs)%2 != 0 {
		panic("makePairs requires key/value pairs")
	}
	out := make([]*api.KVPair, 0, len(kvs)/2)
	for i := 0; i < len(kvs); i += 2 {
		out = append(out, &api.KVPair{Key: kvs[i], Value: []byte(kvs[i+1])})
	}
	return out
}

func TestApply_NoDuplicates_ReturnsAll(t *testing.T) {
	pairs := makePairs("a", "1", "b", "2", "c", "3")
	got := dedupe.Apply(pairs, dedupe.Options{})
	if len(got) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got))
	}
}

func TestApply_EmptyInput_ReturnsEmpty(t *testing.T) {
	got := dedupe.Apply(nil, dedupe.Options{})
	if len(got) != 0 {
		t.Fatalf("expected empty, got %d", len(got))
	}
}

func TestApply_LastWriteWins(t *testing.T) {
	pairs := makePairs("a", "first", "b", "only", "a", "second")
	got := dedupe.Apply(pairs, dedupe.Options{PreferSource: false})
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	for _, p := range got {
		if p.Key == "a" && string(p.Value) != "second" {
			t.Errorf("expected last value 'second', got %q", p.Value)
		}
	}
}

func TestApply_PreferSource_FirstWins(t *testing.T) {
	pairs := makePairs("a", "first", "b", "only", "a", "second")
	got := dedupe.Apply(pairs, dedupe.Options{PreferSource: true})
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	for _, p := range got {
		if p.Key == "a" && string(p.Value) != "first" {
			t.Errorf("expected first value 'first', got %q", p.Value)
		}
	}
}

func TestApply_AllDuplicates_ReturnsSingle(t *testing.T) {
	pairs := makePairs("x", "1", "x", "2", "x", "3")
	got := dedupe.Apply(pairs, dedupe.Options{})
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
}
