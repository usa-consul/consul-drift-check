package diff

import (
	"testing"
)

func makeQuotaResults(keys []string, status string) []Result {
	out := make([]Result, len(keys))
	for i, k := range keys {
		out[i] = Result{Key: k, Status: status}
	}
	return out
}

func TestEvaluateQuota_EmptyResults_ReturnsNil(t *testing.T) {
	got := EvaluateQuota(nil, QuotaOptions{DefaultMax: 1})
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestEvaluateQuota_NoViolation_ReturnsNil(t *testing.T) {
	results := makeQuotaResults([]string{"app/a", "app/b"}, "modified")
	opts := QuotaOptions{Rules: map[string]int{"app": 5}}
	got := EvaluateQuota(results, opts)
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestEvaluateQuota_RuleViolated(t *testing.T) {
	results := makeQuotaResults([]string{"app/a", "app/b", "app/c"}, "modified")
	opts := QuotaOptions{Rules: map[string]int{"app": 2}}
	got := EvaluateQuota(results, opts)
	if len(got) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(got))
	}
	if got[0].Prefix != "app" {
		t.Errorf("unexpected prefix %q", got[0].Prefix)
	}
	if got[0].Allowed != 2 || got[0].Actual != 3 {
		t.Errorf("unexpected counts: %+v", got[0])
	}
}

func TestEvaluateQuota_DefaultMax_Applied(t *testing.T) {
	keys := []string{"svc/a", "svc/b", "svc/c"}
	results := makeQuotaResults(keys, "only_in_source")
	opts := QuotaOptions{DefaultMax: 2}
	got := EvaluateQuota(results, opts)
	if len(got) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(got))
	}
	if got[0].Prefix != "svc" {
		t.Errorf("unexpected prefix %q", got[0].Prefix)
	}
}

func TestEvaluateQuota_ZeroDefaultMax_NoLimit(t *testing.T) {
	keys := []string{"x/1", "x/2", "x/3", "x/4", "x/5"}
	results := makeQuotaResults(keys, "modified")
	opts := QuotaOptions{DefaultMax: 0}
	got := EvaluateQuota(results, opts)
	if got != nil {
		t.Fatalf("expected nil when DefaultMax=0, got %v", got)
	}
}

func TestEvaluateQuota_MultipleViolations_SortedByPrefix(t *testing.T) {
	results := makeQuotaResults(
		[]string{"z/a", "z/b", "z/c", "a/x", "a/y", "a/z"},
		"modified",
	)
	opts := QuotaOptions{DefaultMax: 2}
	got := EvaluateQuota(results, opts)
	if len(got) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(got))
	}
	if got[0].Prefix != "a" || got[1].Prefix != "z" {
		t.Errorf("expected sorted order, got %v %v", got[0].Prefix, got[1].Prefix)
	}
}

func TestQuotaViolation_String(t *testing.T) {
	v := QuotaViolation{Prefix: "app", Allowed: 3, Actual: 7}
	s := v.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}
