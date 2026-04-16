package sanitize_test

import (
	"testing"

	"github.com/hashicorp/consul/api"

	"github.com/myorg/consul-drift-check/internal/sanitize"
)

func makeKVPairs(keys ...string) api.KVPairs {
	pairs := make(api.KVPairs, len(keys))
	for i, k := range keys {
		pairs[i] = &api.KVPair{Key: k, Value: []byte(k)}
	}
	return pairs
}

func TestApply_NoOptions_ReturnsUnchanged(t *testing.T) {
	input := makeKVPairs("app/db/host", "app/db/port")
	out := sanitize.Apply(input, sanitize.Options{})
	if len(out) != 2 || out[0].Key != "app/db/host" {
		t.Fatalf("unexpected result: %v", out)
	}
}

func TestApply_LowerCase(t *testing.T) {
	input := makeKVPairs("App/DB/Host", "APP/DB/PORT")
	out := sanitize.Apply(input, sanitize.Options{LowerCase: true})
	if out[0].Key != "app/db/host" || out[1].Key != "app/db/port" {
		t.Fatalf("keys not lowercased: %v", out)
	}
}

func TestApply_TrimPrefix(t *testing.T) {
	input := makeKVPairs("dc1/app/key", "dc1/app/other")
	out := sanitize.Apply(input, sanitize.Options{TrimPrefix: "dc1"})
	if out[0].Key != "app/key" || out[1].Key != "app/other" {
		t.Fatalf("prefix not trimmed: %v", out)
	}
}

func TestApply_TrimPrefixWithTrailingSlash(t *testing.T) {
	input := makeKVPairs("dc1/app/key")
	out := sanitize.Apply(input, sanitize.Options{TrimPrefix: "dc1/"})
	if out[0].Key != "app/key" {
		t.Fatalf("unexpected key: %s", out[0].Key)
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	input := makeKVPairs("DC1/App/Key")
	sanitize.Apply(input, sanitize.Options{LowerCase: true, TrimPrefix: "dc1"})
	if input[0].Key != "DC1/App/Key" {
		t.Fatal("original pair was mutated")
	}
}

func TestApply_NilEntrySkipped(t *testing.T) {
	pairs := api.KVPairs{nil, {Key: "a/b", Value: []byte("v")}}
	out := sanitize.Apply(pairs, sanitize.Options{})
	if len(out) != 1 || out[0].Key != "a/b" {
		t.Fatalf("nil entry not skipped: %v", out)
	}
}
