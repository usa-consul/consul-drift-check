package normalize_test

import (
	"testing"

	"github.com/hashicorp/consul/api"

	"github.com/someorg/consul-drift-check/internal/normalize"
)

func makePairs(keys ...string) []*api.KVPair {
	out := make([]*api.KVPair, len(keys))
	for i, k := range keys {
		out[i] = &api.KVPair{Key: k, Value: []byte(k + "-val")}
	}
	return out
}

func pairKeys(pairs []*api.KVPair) []string {
	ks := make([]string, len(pairs))
	for i, p := range pairs {
		ks[i] = p.Key
	}
	return ks
}

func TestApply_NoOptions_ReturnsUnchanged(t *testing.T) {
	pairs := makePairs("app/config", "app/secret")
	got := normalize.Apply(pairs, normalize.Options{})
	for i, p := range got {
		if p.Key != pairs[i].Key {
			t.Errorf("expected %q, got %q", pairs[i].Key, p.Key)
		}
	}
}

func TestApply_LowerCase(t *testing.T) {
	pairs := makePairs("App/Config", "APP/SECRET")
	got := normalize.Apply(pairs, normalize.Options{LowerCase: true})
	want := []string{"app/config", "app/secret"}
	for i, k := range pairKeys(got) {
		if k != want[i] {
			t.Errorf("expected %q, got %q", want[i], k)
		}
	}
}

func TestApply_StripPrefix(t *testing.T) {
	pairs := makePairs("dc1/app/config", "dc1/app/secret")
	got := normalize.Apply(pairs, normalize.Options{StripPrefix: "dc1"})
	want := []string{"app/config", "app/secret"}
	for i, k := range pairKeys(got) {
		if k != want[i] {
			t.Errorf("expected %q, got %q", want[i], k)
		}
	}
}

func TestApply_CollapseSlashes(t *testing.T) {
	pairs := makePairs("app//config", "app///secret")
	got := normalize.Apply(pairs, normalize.Options{CollapseSlashes: true})
	want := []string{"app/config", "app/secret"}
	for i, k := range pairKeys(got) {
		if k != want[i] {
			t.Errorf("expected %q, got %q", want[i], k)
		}
	}
}

func TestApply_ValuesUnchanged(t *testing.T) {
	pairs := makePairs("App/Key")
	got := normalize.Apply(pairs, normalize.Options{LowerCase: true})
	if string(got[0].Value) != string(pairs[0].Value) {
		t.Errorf("value should not be modified")
	}
}
