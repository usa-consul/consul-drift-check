package highlight_test

import (
	"testing"

	"github.com/your-org/consul-drift-check/internal/diff"
	"github.com/your-org/consul-drift-check/internal/highlight"
)

func makeResults() []diff.Result {
	return []diff.Result{
		{Key: "app/config", Status: diff.StatusModified},
		{Key: "app/secret", Status: diff.StatusOnlyInSource},
		{Key: "infra/setting", Status: diff.StatusOnlyInDestination},
		{Key: "common/flag", Status: diff.StatusMatch},
	}
}

func TestApply_NoDrift_AllNoneOrModified(t *testing.T) {
	results := makeResults()
	out := highlight.Apply(results, highlight.Options{})
	if len(out) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(out))
	}
	if out[0].Level != highlight.LevelModified {
		t.Errorf("expected modified for status Modified, got %s", out[0].Level)
	}
	if out[3].Level != highlight.LevelNone {
		t.Errorf("expected none for status Match, got %s", out[3].Level)
	}
}

func TestApply_CriticalPrefix_OverridesLevel(t *testing.T) {
	results := []diff.Result{
		{Key: "infra/setting", Status: diff.StatusMatch},
	}
	out := highlight.Apply(results, highlight.Options{
		CriticalPrefixes: []string{"infra/"},
	})
	if out[0].Level != highlight.LevelCritical {
		t.Errorf("expected critical, got %s", out[0].Level)
	}
}

func TestApply_ModifiedOnly_SkipsSourceDest(t *testing.T) {
	results := []diff.Result{
		{Key: "app/secret", Status: diff.StatusOnlyInSource},
		{Key: "app/config", Status: diff.StatusModified},
	}
	out := highlight.Apply(results, highlight.Options{ModifiedOnly: true})
	if out[0].Level != highlight.LevelNone {
		t.Errorf("expected none for OnlyInSource with ModifiedOnly, got %s", out[0].Level)
	}
	if out[1].Level != highlight.LevelModified {
		t.Errorf("expected modified, got %s", out[1].Level)
	}
}

func TestApply_EmptyResults_ReturnsEmpty(t *testing.T) {
	out := highlight.Apply(nil, highlight.Options{})
	if len(out) != 0 {
		t.Errorf("expected empty slice, got %d elements", len(out))
	}
}
