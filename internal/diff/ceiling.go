package diff

import "sort"

// CeilingOptions controls how the ceiling is applied to drift results.
type CeilingOptions struct {
	// MaxPerPrefix is the maximum number of results allowed per top-level prefix.
	// Zero means no ceiling is applied.
	MaxPerPrefix int
	// OrderByKey sorts results within each prefix before applying the ceiling.
	OrderByKey bool
}

// Ceiling limits the number of drift results per top-level KV prefix.
// Results that exceed the per-prefix cap are dropped. The relative order of
// kept results is preserved unless OrderByKey is set.
func Ceiling(results []Result, opts CeilingOptions) []Result {
	if len(results) == 0 || opts.MaxPerPrefix <= 0 {
		return results
	}

	working := make([]Result, len(results))
	copy(working, results)

	if opts.OrderByKey {
		sort.Slice(working, func(i, j int) bool {
			return working[i].Key < working[j].Key
		})
	}

	counts := make(map[string]int)
	out := make([]Result, 0, len(working))
	for _, r := range working {
		seg := topCeilingSegment(r.Key)
		if counts[seg] < opts.MaxPerPrefix {
			out = append(out, r)
			counts[seg]++
		}
	}
	return out
}

func topCeilingSegment(key string) string {
	if len(key) > 0 && key[0] == '/' {
		key = key[1:]
	}
	for i, c := range key {
		if c == '/' {
			return key[:i]
		}
	}
	return key
}
