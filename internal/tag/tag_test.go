package tag_test

import (
	"testing"

	"github.com/your-org/consul-drift-check/internal/diff"
	"github.com/your-org/consul-drift-check/internal/tag"
)

func makeResults(keys ...string) []diff.Result {
	out := make([]diff.Result, len(keys))
	for i, k := range keys {
		out[i] = diff.Result{Key: k, Status: diff.Modified}
	}
	return out
}

func TestApply_NoRules_ReturnsEmptyTags(t *testing.T) {
	results := makeResults("app/config", "db/host")
	tagged := tag.Apply(results, tag.Options{})
	for _, tr := range tagged {
		if len(tr.Tags) != 0 {
			t.Errorf("expected no tags for %q, got %v", tr.Key, tr.Tags)
		}
	}
}

func TestApply_MatchingPrefix_AssignsTags(t *testing.T) {
	opts := tag.Options{
		Rules: []tag.Rule{
			{Prefix: "app/", Tags: []string{"team:backend", "env:prod"}},
		},
	}
	tagged := tag.Apply(makeResults("app/config"), opts)
	if len(tagged) != 1 {
		t.Fatalf("expected 1 result")
	}
	if len(tagged[0].Tags) != 2 || tagged[0].Tags[0] != "team:backend" {
		t.Errorf("unexpected tags: %v", tagged[0].Tags)
	}
}

func TestApply_FirstRuleWins(t *testing.T) {
	opts := tag.Options{
		Rules: []tag.Rule{
			{Prefix: "app/", Tags: []string{"first"}},
			{Prefix: "app/config", Tags: []string{"second"}},
		},
	}
	tagged := tag.Apply(makeResults("app/config"), opts)
	if tagged[0].Tags[0] != "first" {
		t.Errorf("expected first rule to win, got %v", tagged[0].Tags)
	}
}

func TestApply_NoMatchingPrefix_EmptyTags(t *testing.T) {
	opts := tag.Options{
		Rules: []tag.Rule{
			{Prefix: "db/", Tags: []string{"team:data"}},
		},
	}
	tagged := tag.Apply(makeResults("app/config"), opts)
	if len(tagged[0].Tags) != 0 {
		t.Errorf("expected empty tags, got %v", tagged[0].Tags)
	}
}

func TestApply_TagsAreCopied(t *testing.T) {
	rule := tag.Rule{Prefix: "app/", Tags: []string{"original"}}
	opts := tag.Options{Rules: []tag.Rule{rule}}
	tagged := tag.Apply(makeResults("app/x"), opts)
	tagged[0].Tags[0] = "mutated"
	if opts.Rules[0].Tags[0] != "original" {
		t.Error("Apply mutated the original rule tags")
	}
}
