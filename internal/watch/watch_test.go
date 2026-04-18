package watch_test

import (
	"testing"

	"github.com/hashicorp/consul/api"
	"github.com/nicholasgasior/consul-drift-check/internal/watch"
)

func makePairs(data map[string]string) api.KVPairs {
	pairs := make(api.KVPairs, 0, len(data))
	for k, v := range data {
		pairs = append(pairs, &api.KVPair{Key: k, Value: []byte(v)})
	}
	return pairs
}

func TestDiff_NoChanges(t *testing.T) {
	prev := makePairs(map[string]string{"a": "1", "b": "2"})
	curr := makePairs(map[string]string{"a": "1", "b": "2"})
	changes := watch.Diff(prev, curr)
	if len(changes) != 0 {
		t.Fatalf("expected 0 changes, got %d", len(changes))
	}
}

func TestDiff_DetectsAdded(t *testing.T) {
	prev := makePairs(map[string]string{"a": "1"})
	curr := makePairs(map[string]string{"a": "1", "b": "2"})
	changes := watch.Diff(prev, curr)
	if len(changes) != 1 || changes[0].Kind != watch.Added || changes[0].Key != "b" {
		t.Fatalf("expected added 'b', got %+v", changes)
	}
}

func TestDiff_DetectsRemoved(t *testing.T) {
	prev := makePairs(map[string]string{"a": "1", "b": "2"})
	curr := makePairs(map[string]string{"a": "1"})
	changes := watch.Diff(prev, curr)
	if len(changes) != 1 || changes[0].Kind != watch.Removed || changes[0].Key != "b" {
		t.Fatalf("expected removed 'b', got %+v", changes)
	}
}

func TestDiff_DetectsModified(t *testing.T) {
	prev := makePairs(map[string]string{"a": "old"})
	curr := makePairs(map[string]string{"a": "new"})
	changes := watch.Diff(prev, curr)
	if len(changes) != 1 || changes[0].Kind != watch.Modified {
		t.Fatalf("expected modified 'a', got %+v", changes)
	}
	if string(changes[0].OldValue) != "old" || string(changes[0].NewValue) != "new" {
		t.Fatalf("unexpected values: %+v", changes[0])
	}
}

func TestDiff_EmptyInputs(t *testing.T) {
	changes := watch.Diff(nil, nil)
	if len(changes) != 0 {
		t.Fatalf("expected 0 changes, got %d", len(changes))
	}
}

func TestDiff_NilPairsSkipped(t *testing.T) {
	prev := api.KVPairs{nil, {Key: "a", Value: []byte("1")}}
	curr := api.KVPairs{{Key: "a", Value: []byte("1")}}
	changes := watch.Diff(prev, curr)
	if len(changes) != 0 {
		t.Fatalf("expected 0 changes, got %d", len(changes))
	}
}
