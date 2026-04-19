package diff

import (
	"testing"

	"github.com/hashicorp/consul/api"
)

func makeKVPairs(kvs map[string]string) []*api.KVPair {
	pairs := make([]*api.KVPair, 0, len(kvs))
	for k, v := range kvs {
		pairs = append(pairs, &api.KVPair{Key: k, Value: []byte(v)})
	}
	return pairs
}

func TestDigest_EmptyPairs(t *testing.T) {
	e := Digest("app/", nil)
	if e.Digest == "" {
		t.Fatal("expected non-empty digest")
	}
	if e.KeyCount != 0 {
		t.Fatalf("expected 0 keys, got %d", e.KeyCount)
	}
}

func TestDigest_SameContentSameDigest(t *testing.T) {
	pairs := makeKVPairs(map[string]string{"a": "1", "b": "2"})
	a := Digest("app/", pairs)
	b := Digest("app/", pairs)
	if a.Digest != b.Digest {
		t.Fatal("identical pairs produced different digests")
	}
}

func TestDigest_OrderIndependent(t *testing.T) {
	p1 := []*api.KVPair{{Key: "a", Value: []byte("1")}, {Key: "b", Value: []byte("2")}}
	p2 := []*api.KVPair{{Key: "b", Value: []byte("2")}, {Key: "a", Value: []byte("1")}}
	if Digest("x", p1).Digest != Digest("x", p2).Digest {
		t.Fatal("digest should be order-independent")
	}
}

func TestDigest_DifferentValuesDifferentDigest(t *testing.T) {
	a := Digest("p", makeKVPairs(map[string]string{"k": "v1"}))
	b := Digest("p", makeKVPairs(map[string]string{"k": "v2"}))
	if a.Digest == b.Digest {
		t.Fatal("different values should produce different digests")
	}
}

func TestDigestsEqual_TrueForSame(t *testing.T) {
	pairs := makeKVPairs(map[string]string{"x": "y"})
	if !DigestsEqual(Digest("p", pairs), Digest("p", pairs)) {
		t.Fatal("expected equal digests")
	}
}

func TestDigestsEqual_FalseForDifferent(t *testing.T) {
	a := Digest("p", makeKVPairs(map[string]string{"k": "1"}))
	b := Digest("p", makeKVPairs(map[string]string{"k": "2"}))
	if DigestsEqual(a, b) {
		t.Fatal("expected unequal digests")
	}
}
