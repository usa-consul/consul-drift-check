package summarize_test

import (
	"fmt"
	"testing"

	"github.com/your-org/consul-drift-check/internal/diff"
	"github.com/your-org/consul-drift-check/internal/summarize"
)

func BenchmarkBuild_LargeResultSet(b *testing.B) {
	statuses := []diff.Status{
		diff.StatusOnlyInSource,
		diff.StatusOnlyInDestination,
		diff.StatusModified,
	}
	results := make([]diff.Result, 1000)
	for i := range results {
		results[i] = diff.Result{
			Key:    fmt.Sprintf("service/config/key-%d", i),
			Status: statuses[i%3],
		}
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		summarize.Build(results)
	}
}
