package diff

import (
	"testing"

	"github.com/hashicorp/consul/api"
)

func makePruneKV(key string, value []byte) *api.KVPair {
	return &api.KVPair{Key: key, Value: value}
}

func pruneMap(pairs ...*api.KVPair) map[string]*api.KVPair {
	m := make(map[string]*api.KVPair, len(pairs))
	for _, p := range pairs {
		if p != nil {
			m[p.Key] = p
		}
	}
	return m
}

func TestPrune_EmptyInput_ReturnsNil(t *testing.T) {
	out := Prune(nil, nil, nil, PruneOptions{})
	if out != nil {
		t.Fatalf("expected nil, got %v", out)
	}
}

func TestPrune_NoOptions_ReturnsAll(t *testing.T) {
	a := makePruneKV("app/foo", []byte("v1"))
	b := makePruneKV("app/bar", []byte("v2"))
	src := pruneMap(a, b)
	dst := pruneMap()
	out := Prune([]*api.KVPair{a, b}, src, dst, PruneOptions{})
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}

func TestPrune_DropEmptyValues_RemovesBothEmpty(t *testing.T) {
	a := makePruneKV("app/foo", []byte{})
	b := makePruneKV("app/bar", []byte("val"))
	src := pruneMap(a, b)
	dst := pruneMap(a)
	out := Prune([]*api.KVPair{a, b}, src, dst, PruneOptions{DropEmptyValues: true})
	if len(out) != 1 || out[0].Key != "app/bar" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestPrune_OnlyStatus_FiltersCorrectly(t *testing.T) {
	a := makePruneKV("app/only-src", []byte("x"))
	b := makePruneKV("app/only-dst", []byte("y"))
	src := pruneMap(a)
	dst := pruneMap(b)
	out := Prune([]*api.KVPair{a, b}, src, dst, PruneOptions{OnlyStatus: []string{"source_only"}})
	if len(out) != 1 || out[0].Key != "app/only-src" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestPrune_StalePrefixes_DropsMatching(t *testing.T) {
	a := makePruneKV("stale/foo", []byte("x"))
	b := makePruneKV("live/bar", []byte("y"))
	src := pruneMap(a, b)
	dst := pruneMap()
	out := Prune([]*api.KVPair{a, b}, src, dst, PruneOptions{StalePrefixes: []string{"stale/"}})
	if len(out) != 1 || out[0].Key != "live/bar" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestPrune_NilEntry_Skipped(t *testing.T) {
	a := makePruneKV("app/foo", []byte("v"))
	src := pruneMap(a)
	out := Prune([]*api.KVPair{nil, a}, src, pruneMap(), PruneOptions{})
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
}
