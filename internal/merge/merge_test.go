package merge_test

import (
	"testing"

	"github.com/hashicorp/consul/api"

	"github.com/your-org/consul-drift-check/internal/merge"
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

func TestApply_NoConflict_ReturnsMerged(t *testing.T) {
	src := makePairs("a", "1")
	dst := makePairs("b", "2")
	got := merge.Apply(src, dst, merge.Options{})
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}

func TestApply_SourceWins_OnConflict(t *testing.T) {
	src := makePairs("key", "from-src")
	dst := makePairs("key", "from-dst")
	got := merge.Apply(src, dst, merge.Options{Strategy: merge.StrategySourceWins})
	if string(got[0].Value) != "from-src" {
		t.Errorf("expected from-src, got %s", got[0].Value)
	}
}

func TestApply_DestWins_OnConflict(t *testing.T) {
	src := makePairs("key", "from-src")
	dst := makePairs("key", "from-dst")
	got := merge.Apply(src, dst, merge.Options{Strategy: merge.StrategyDestWins})
	if string(got[0].Value) != "from-dst" {
		t.Errorf("expected from-dst, got %s", got[0].Value)
	}
}

func TestApply_Latest_PicksHigherModifyIndex(t *testing.T) {
	src := []*api.KVPair{{Key: "key", Value: []byte("old"), ModifyIndex: 1}}
	dst := []*api.KVPair{{Key: "key", Value: []byte("new"), ModifyIndex: 5}}
	got := merge.Apply(src, dst, merge.Options{Strategy: merge.StrategyLatest})
	if string(got[0].Value) != "new" {
		t.Errorf("expected new, got %s", got[0].Value)
	}
}

func TestApply_NilEntries_Skipped(t *testing.T) {
	src := []*api.KVPair{nil, {Key: "a", Value: []byte("1")}}
	dst := []*api.KVPair{{Key: "b", Value: []byte("2")}, nil}
	got := merge.Apply(src, dst, merge.Options{})
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestApply_ResultIsSortedByKey(t *testing.T) {
	src := makePairs("z", "1", "a", "2", "m", "3")
	got := merge.Apply(src, nil, merge.Options{})
	for i := 1; i < len(got); i++ {
		if got[i].Key < got[i-1].Key {
			t.Errorf("result not sorted at index %d: %s < %s", i, got[i].Key, got[i-1].Key)
		}
	}
}
