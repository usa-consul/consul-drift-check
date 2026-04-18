package annotate_test

import (
	"testing"

	"github.com/example/consul-drift-check/internal/annotate"
	"github.com/example/consul-drift-check/internal/diff"
)

func makeResults(keys ...string) []diff.Result {
	out := make([]diff.Result, len(keys))
	for i, k := range keys {
		out[i] = diff.Result{Key: k, Status: diff.Modified}
	}
	return out
}

func TestApply_NoRules_ReturnsEmptyAnnotations(t *testing.T) {
	res := annotate.Apply(makeResults("app/config"), nil)
	if len(res) != 1 {
		t.Fatalf("expected 1 result, got %d", len(res))
	}
	if len(res[0].Annotations) != 0 {
		t.Errorf("expected empty annotations, got %v", res[0].Annotations)
	}
}

func TestApply_MatchingPrefix_AssignsAnnotations(t *testing.T) {
	rules := []annotate.Rule{
		{Prefix: "app/", Annotations: map[string]string{"team": "platform"}},
	}
	res := annotate.Apply(makeResults("app/config"), rules)
	if res[0].Annotations["team"] != "platform" {
		t.Errorf("expected team=platform, got %v", res[0].Annotations)
	}
}

func TestApply_FirstRuleWins(t *testing.T) {
	rules := []annotate.Rule{
		{Prefix: "app/", Annotations: map[string]string{"owner": "alpha"}},
		{Prefix: "app/config", Annotations: map[string]string{"owner": "beta"}},
	}
	res := annotate.Apply(makeResults("app/config"), rules)
	if res[0].Annotations["owner"] != "alpha" {
		t.Errorf("expected owner=alpha, got %s", res[0].Annotations["owner"])
	}
}

func TestApply_NoMatchingPrefix_EmptyAnnotations(t *testing.T) {
	rules := []annotate.Rule{
		{Prefix: "db/", Annotations: map[string]string{"tier": "data"}},
	}
	res := annotate.Apply(makeResults("app/config"), rules)
	if len(res[0].Annotations) != 0 {
		t.Errorf("expected empty annotations, got %v", res[0].Annotations)
	}
}

func TestApply_AnnotationsAreCopied(t *testing.T) {
	orig := map[string]string{"env": "prod"}
	rules := []annotate.Rule{{Prefix: "app/", Annotations: orig}}
	res := annotate.Apply(makeResults("app/x", "app/y"), rules)
	res[0].Annotations["env"] = "staging"
	if res[1].Annotations["env"] != "prod" {
		t.Errorf("annotations should be independent copies")
	}
}
