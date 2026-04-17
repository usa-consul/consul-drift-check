package flatten_test

import (
	"testing"

	"github.com/hashicorp/consul/api"

	"github.com/example/consul-drift-check/internal/flatten"
)

func makePairs(kvs ...string) api.KVPairs {
	if len(kvs)%2 != 0 {
		panic("makePairs requires even number of arguments")
	}
	pairs := make(api.KVPairs, 0, len(kvs)/2)
	for i := 0; i < len(kvs); i += 2 {
		pairs = append(pairs, &api.KVPair{Key: kvs[i], Value: []byte(kvs[i+1])})
	}
	return pairs
}

func TestApply_NoOptions_ReturnsFlatMap(t *testing.T) {
	pairs := makePairs("app/db/host", "localhost", "app/db/port", "5432")
	out := flatten.Apply(pairs, flatten.Options{})
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if string(out["app/db/host"]) != "localhost" {
		t.Errorf("unexpected value for app/db/host")
	}
}

func TestApply_StripPrefix_RemovesLeadingSegment(t *testing.T) {
	pairs := makePairs("prod/app/key", "val")
	out := flatten.Apply(pairs, flatten.Options{StripPrefix: "prod"})
	if _, ok := out["app/key"]; !ok {
		t.Errorf("expected key 'app/key' after stripping prefix, got %v", out)
	}
}

func TestApply_NilValue_StoredAsEmpty(t *testing.T) {
	pairs := api.KVPairs{{Key: "some/key", Value: nil}}
	out := flatten.Apply(pairs, flatten.Options{})
	v, ok := out["some/key"]
	if !ok {
		t.Fatal("expected key to be present")
	}
	if len(v) != 0 {
		t.Errorf("expected empty value, got %v", v)
	}
}

func TestApply_NilPair_Skipped(t *testing.T) {
	pairs := api.KVPairs{nil, {Key: "a/b", Value: []byte("x")}}
	out := flatten.Apply(pairs, flatten.Options{})
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
}

func TestApply_EmptyInput_ReturnsEmptyMap(t *testing.T) {
	out := flatten.Apply(api.KVPairs{}, flatten.Options{})
	if len(out) != 0 {
		t.Errorf("expected empty map")
	}
}

func TestApply_PrefixOnlyKey_Skipped(t *testing.T) {
	pairs := makePairs("prod/", "")
	out := flatten.Apply(pairs, flatten.Options{StripPrefix: "prod"})
	if len(out) != 0 {
		t.Errorf("expected empty map after stripping prefix-only key, got %v", out)
	}
}
