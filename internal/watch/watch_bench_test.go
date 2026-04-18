package watch_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/consul/api"
	"github.com/nicholasgasior/consul-drift-check/internal/watch"
)

func BenchmarkDiff_LargePairSet(b *testing.B) {
	const n = 5000
	prev := make(api.KVPairs, n)
	curr := make(api.KVPairs, n)
	for i := 0; i < n; i++ {
		key := fmt.Sprintf("config/key/%d", i)
		prev[i] = &api.KVPair{Key: key, Value: []byte("value-old")}
		val := "value-old"
		if i%10 == 0 {
			val = "value-new"
		}
		curr[i] = &api.KVPair{Key: key, Value: []byte(val)}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = watch.Diff(prev, curr)
	}
}
