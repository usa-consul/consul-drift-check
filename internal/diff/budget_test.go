package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeBudgetResults(statuses ...string) []Result {
	out := make([]Result, len(statuses))
	for i, s := range statuses {
		out[i] = Result{Key: "k", Status: s}
	}
	return out
}

func TestEvaluateBudget_NoDrift_NotExceeded(t *testing.T) {
	res := EvaluateBudget(makeBudgetResults(), BudgetOptions{MaxDriftCount: 5})
	assert.False(t, res.Exceeded)
	assert.Equal(t, 0, res.TotalDrift)
	assert.Empty(t, res.Violations)
}

func TestEvaluateBudget_UnderLimit_NotExceeded(t *testing.T) {
	results := makeBudgetResults("modified", "source-only")
	res := EvaluateBudget(results, BudgetOptions{MaxDriftCount: 5})
	assert.False(t, res.Exceeded)
	assert.Equal(t, 2, res.TotalDrift)
}

func TestEvaluateBudget_TotalExceeded(t *testing.T) {
	results := makeBudgetResults("modified", "modified", "source-only")
	res := EvaluateBudget(results, BudgetOptions{MaxDriftCount: 2})
	assert.True(t, res.Exceeded)
	assert.Contains(t, res.Violations, "total drift count exceeded")
}

func TestEvaluateBudget_ModifiedExceeded(t *testing.T) {
	results := makeBudgetResults("modified", "modified", "modified")
	res := EvaluateBudget(results, BudgetOptions{MaxModified: 2})
	assert.True(t, res.Exceeded)
	assert.Contains(t, res.Violations, "modified key count exceeded")
}

func TestEvaluateBudget_MissingExceeded(t *testing.T) {
	results := makeBudgetResults("source-only", "destination-only", "source-only")
	res := EvaluateBudget(results, BudgetOptions{MaxMissing: 1})
	assert.True(t, res.Exceeded)
	assert.Contains(t, res.Violations, "missing key count exceeded")
	assert.Equal(t, 3, res.Missing)
}

func TestEvaluateBudget_MultipleViolations(t *testing.T) {
	results := makeBudgetResults("modified", "modified", "source-only", "source-only")
	res := EvaluateBudget(results, BudgetOptions{
		MaxDriftCount: 1,
		MaxModified:   1,
		MaxMissing:    1,
	})
	assert.True(t, res.Exceeded)
	assert.Len(t, res.Violations, 3)
}

func TestEvaluateBudget_ZeroLimits_NeverExceeded(t *testing.T) {
	results := makeBudgetResults("modified", "modified", "source-only")
	res := EvaluateBudget(results, BudgetOptions{})
	assert.False(t, res.Exceeded)
	assert.Empty(t, res.Violations)
}
