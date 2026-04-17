package classify_test

import (
	"testing"

	"github.com/your-org/consul-drift-check/internal/classify"
	"github.com/your-org/consul-drift-check/internal/diff"
)

func makeResults(keys ...string) []diff.Result {
	out := make([]diff.Result, len(keys))
	for i, k := range keys {
		out[i] = diff.Result{Key: k, Status: diff.StatusModified}
	}
	return out
}

func TestApply_NoRules_DefaultsToInfo(t *testing.T) {
	res := classify.Apply(makeResults("app/config"), nil)
	if res[0].Level != classify.LevelInfo {
		t.Fatalf("expected info, got %s", res[0].Level)
	}
}

func TestApply_MatchingPrefix_AssignsLevel(t *testing.T) {
	rules := []classify.Rule{
		{Prefix: "prod/secrets/", Level: classify.LevelCritical},
	}
	res := classify.Apply(makeResults("prod/secrets/db"), rules)
	if res[0].Level != classify.LevelCritical {
		t.Fatalf("expected critical, got %s", res[0].Level)
	}
}

func TestApply_FirstRuleWins(t *testing.T) {
	rules := []classify.Rule{
		{Prefix: "prod/", Level: classify.LevelWarning},
		{Prefix: "prod/secrets/", Level: classify.LevelCritical},
	}
	res := classify.Apply(makeResults("prod/secrets/key"), rules)
	if res[0].Level != classify.LevelWarning {
		t.Fatalf("expected warning (first rule wins), got %s", res[0].Level)
	}
}

func TestApply_MixedKeys(t *testing.T) {
	rules := []classify.Rule{
		{Prefix: "prod/", Level: classify.LevelCritical},
		{Prefix: "staging/", Level: classify.LevelWarning},
	}
	input := makeResults("prod/cfg", "staging/cfg", "dev/cfg")
	res := classify.Apply(input, rules)
	expected := []classify.Level{classify.LevelCritical, classify.LevelWarning, classify.LevelInfo}
	for i, r := range res {
		if r.Level != expected[i] {
			t.Errorf("[%d] expected %s, got %s", i, expected[i], r.Level)
		}
	}
}

func TestApply_EmptyResults(t *testing.T) {
	res := classify.Apply(nil, []classify.Rule{{Prefix: "x/", Level: classify.LevelCritical}})
	if len(res) != 0 {
		t.Fatalf("expected empty slice")
	}
}

func TestApply_ExactPrefixMatch(t *testing.T) {
	// Ensure a rule only matches keys that start with the prefix,
	// not keys where the prefix appears elsewhere in the path.
	rules := []classify.Rule{
		{Prefix: "prod/", Level: classify.LevelCritical},
	}
	input := makeResults("notprod/cfg", "staging/prod/cfg")
	res := classify.Apply(input, rules)
	for i, r := range res {
		if r.Level != classify.LevelInfo {
			t.Errorf("[%d] expected info for non-prefix match, got %s", i, r.Level)
		}
	}
}
