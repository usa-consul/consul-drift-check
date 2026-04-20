package diff

import (
	"testing"

	"github.com/example/consul-drift-check/internal/consul"
)

func makeNormalizeResults(kvs ...string) []Result {
	results := make([]Result, 0, len(kvs))
	for _, k := range kvs {
		results = append(results, Result{Key: k, Status: "modified"})
	}
	return results
}

func TestNormalizeResults_NoOptions_ReturnsUnchanged(t *testing.T) {
	input := makeNormalizeResults("app/db", "app/cache")
	out := NormalizeResults(input, NormalizeOptions{})
	if len(out) != 2 || out[0].Key != "app/db" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestNormalizeResults_StripPrefix(t *testing.T) {
	input := makeNormalizeResults("prod/app/db", "prod/app/cache")
	out := NormalizeResults(input, NormalizeOptions{StripPrefix: "prod"})
	if out[0].Key != "app/db" || out[1].Key != "app/cache" {
		t.Fatalf("unexpected keys: %v", out)
	}
}

func TestNormalizeResults_LowerCase(t *testing.T) {
	input := makeNormalizeResults("App/DB", "App/Cache")
	out := NormalizeResults(input, NormalizeOptions{LowerCase: true})
	if out[0].Key != "app/db" || out[1].Key != "app/cache" {
		t.Fatalf("unexpected keys: %v", out)
	}
}

func TestNormalizeResults_EmptyInput_ReturnsNil(t *testing.T) {
	out := NormalizeResults(nil, NormalizeOptions{StripPrefix: "x"})
	if out != nil {
		t.Fatalf("expected nil, got %v", out)
	}
}

func TestNormalizePairs_StripPrefix(t *testing.T) {
	pairs := []*consul.KVPair{
		{Key: "staging/svc/port", Value: []byte("8080")},
		{Key: "staging/svc/host", Value: []byte("localhost")},
	}
	out := NormalizePairs(pairs, NormalizeOptions{StripPrefix: "staging"})
	if len(out) != 2 || out[0].Key != "svc/port" {
		t.Fatalf("unexpected pairs: %v", out)
	}
}

func TestNormalizePairs_NilPair_Skipped(t *testing.T) {
	pairs := []*consul.KVPair{nil, {Key: "a/b", Value: []byte("v")}}
	out := NormalizePairs(pairs, NormalizeOptions{})
	if len(out) != 1 || out[0].Key != "a/b" {
		t.Fatalf("expected 1 pair, got %v", out)
	}
}

func TestNormalizePairs_EmptyInput_ReturnsNil(t *testing.T) {
	out := NormalizePairs(nil, NormalizeOptions{})
	if out != nil {
		t.Fatalf("expected nil, got %v", out)
	}
}
