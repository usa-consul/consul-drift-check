package diff

import (
	"fmt"
	"testing"
)

func BenchmarkCeiling_LargeResultSet(b *testing.B) {
	const total = 10_000
	results := make([]Result, total)
	prefixes := []string{"app", "db", "cache", "queue", "infra"}
	for i := 0; i < total; i++ {
		pfx := prefixes[i%len(prefixes)]
		results[i] = Result{
			Key:    fmt.Sprintf("%s/key-%d", pfx, i),
			Status: "modified",
		}
	}
	opts := CeilingOptions{MaxPerPrefix: 100, OrderByKey: true}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Ceiling(results, opts)
	}
}
