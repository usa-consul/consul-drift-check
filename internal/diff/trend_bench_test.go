package diff

import (
	"testing"
	"time"
)

func BenchmarkBuildTrend_LargeHistory(b *testing.B) {
	now := time.Now().UTC()
	results := make([]Result, 200)
	for i := range results {
		results[i] = Result{Key: "key", Status: "modified"}
	}
	entries := make([]HistoryEntry, 100)
	for i := range entries {
		entries[i] = HistoryEntry{
			Timestamp: now.Add(time.Duration(i) * time.Minute),
			Results:   results,
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BuildTrend(entries)
	}
}
