package diff

import (
	"fmt"
	"testing"
)

func BenchmarkEvaluateQuota_LargeResultSet(b *testing.B) {
	const n = 10_000
	results := make([]Result, n)
	for i := range results {
		results[i] = Result{
			Key:    fmt.Sprintf("prefix%d/key%d", i%20, i),
			Status: "modified",
		}
	}
	opts := QuotaOptions{
		Rules:      map[string]int{"prefix0": 100, "prefix1": 200},
		DefaultMax: 500,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EvaluateQuota(results, opts)
	}
}
