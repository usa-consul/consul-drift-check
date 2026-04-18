package diff

import (
	"sort"
	"strings"
)

// RollupEntry summarises drift results aggregated under a common key prefix.
type RollupEntry struct {
	Prefix   string
	Added    int
	Removed  int
	Modified int
	Total    int
}

// Rollup aggregates a slice of Result values by their top-level key segment.
// Each entry in the returned slice represents one distinct prefix bucket.
func Rollup(results []Result) []RollupEntry {
	counts := make(map[string]*RollupEntry)

	for _, r := range results {
		prefix := topSegment(r.Key)
		e, ok := counts[prefix]
		if !ok {
			e = &RollupEntry{Prefix: prefix}
			counts[prefix] = e
		}
		switch r.Status {
		case StatusOnlyInSource:
			e.Added++
		case StatusOnlyInDestination:
			e.Removed++
		case StatusModified:
			e.Modified++
		}
		e.Total++
	}

	out := make([]RollupEntry, 0, len(counts))
	for _, e := range counts {
		out = append(out, *e)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Prefix < out[j].Prefix
	})
	return out
}

func topSegment(key string) string {
	key = strings.TrimPrefix(key, "/")
	if idx := strings.Index(key, "/"); idx >= 0 {
		return key[:idx]
	}
	return key
}
