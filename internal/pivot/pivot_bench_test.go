package pivot_test

import (
	"fmt"
	"testing"

	"github.com/your-org/consul-drift-check/internal/diff"
	"github.com/your-org/consul-drift-check/internal/pivot"
)

func BenchmarkBuild_LargeMultiDC(b *testing.B) {
	const keysPerDC = 500
	const dcCount = 6

	dcResults := make(map[string][]diff.Result, dcCount)
	for d := 0; d < dcCount; d++ {
		label := fmt.Sprintf("dc%d", d+1)
		results := make([]diff.Result, keysPerDC)
		for k := 0; k < keysPerDC; k++ {
			results[k] = diff.Result{
				Key:              fmt.Sprintf("service/%d/config", k),
				Status:           diff.Modified,
				DestinationValue: fmt.Sprintf("val-%d-%d", d, k),
			}
		}
		dcResults[label] = results
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pivot.Build(dcResults)
	}
}
