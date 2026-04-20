package diff

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkDetectStale_LargeResultSet(b *testing.B) {
	const n = 10_000
	now := time.Now()

	results := make([]Result, n)
	seenAt := make(map[string]time.Time, n)
	for i := 0; i < n; i++ {
		key := fmt.Sprintf("app/service-%d/config", i)
		results[i] = Result{Key: key, Status: "modified"}
		seenAt[key] = now.Add(-time.Duration(i) * time.Minute)
	}

	opts := StaleOptions{
		MaxAge:          12 * time.Hour,
		ExcludePrefixes: []string{"app/service-9"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = DetectStale(results, seenAt, now, opts)
	}
}
