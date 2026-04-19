package diff

import (
	"testing"

	"github.com/hashicorp/consul/api"
)

func makeExcludeResults(keyStatus ...string) []Result {
	if len(keyStatus)%2 != 0 {
		panic("makeExcludeResults: need key,status pairs")
	}
	out := make([]Result, 0, len(keyStatus)/2)
	for i := 0; i < len(keyStatus); i += 2 {
		out = append(out, Result{Key: keyStatus[i], Status: keyStatus[i+1]})
	}
	return out
}

func TestExclude_NoOptions_ReturnsAll(t *testing.T) {
	results := makeExcludeResults("a", "modified", "b", "only_in_source")
	out := Exclude(results, ExcludeOptions{})
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}

func TestExclude_ByStatus_DropsMatching(t *testing.T) {
	results := makeExcludeResults("a", "modified", "b", "only_in_source", "c", "only_in_dest")
	out := Exclude(results, ExcludeOptions{Statuses: []string{"only_in_source", "only_in_dest"}})
	if len(out) != 1 || out[0].Key != "a" {
		t.Fatalf("unexpected results: %+v", out)
	}
}

func TestExclude_ByPrefix_DropsMatching(t *testing.T) {
	results := makeExcludeResults("app/cfg", "modified", "sys/cfg", "modified", "app/secret", "only_in_source")
	out := Exclude(results, ExcludeOptions{Prefixes: []string{"app/"}})
	if len(out) != 1 || out[0].Key != "sys/cfg" {
		t.Fatalf("unexpected results: %+v", out)
	}
}

func TestExclude_StatusCaseInsensitive(t *testing.T) {
	results := makeExcludeResults("k", "Modified")
	out := Exclude(results, ExcludeOptions{Statuses: []string{"modified"}})
	if len(out) != 0 {
		t.Fatalf("expected empty, got %+v", out)
	}
}

func TestExclude_EmptyResults_ReturnsEmpty(t *testing.T) {
	out := Exclude(nil, ExcludeOptions{Statuses: []string{"modified"}})
	if len(out) != 0 {
		t.Fatalf("expected empty slice")
	}
}

func TestExcludeFromPairs_NoPrefixes_ReturnsAll(t *testing.T) {
	pairs := []*api.KVPair{{Key: "a"}, {Key: "b"}}
	out := ExcludeFromPairs(pairs, nil)
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}

func TestExcludeFromPairs_MatchingPrefix_Removed(t *testing.T) {
	pairs := []*api.KVPair{{Key: "secret/db"}, {Key: "config/db"}}
	out := ExcludeFromPairs(pairs, []string{"secret/"})
	if len(out) != 1 || out[0].Key != "config/db" {
		t.Fatalf("unexpected pairs: %+v", out)
	}
}
