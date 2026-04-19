package diff

import (
	"fmt"
	"testing"
)

func BenchmarkEvaluateBudget_LargeResultSet(b *testing.B) {
	statuses := []string{"modified", "source-only", "destination-only", "match"}
	results := make([]Result, 10_000)
	for i := range results {
		results[i] = Result{
			Key:    fmt.Sprintf("prefix/key-%d", i),
			Status: statuses[i%len(statuses)],
		}
	}
	opts := BudgetOptions{
		MaxDriftCount: 5000,
		MaxModified:   2500,
		MaxMissing:    2500,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EvaluateBudget(results, opts)
	}
}
