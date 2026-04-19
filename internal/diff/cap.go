package diff

import "sort"

// CapOptions controls how Cap trims a result set.
type CapOptions struct {
	// MaxResults is the maximum number of results to return.
	// Zero means no limit.
	MaxResults int
	// OrderByKey sorts alphabetically before capping when true.
	// Otherwise the original order is preserved.
	OrderByKey bool
}

// Cap returns at most opts.MaxResults entries from results.
// When OrderByKey is true the slice is sorted by key before trimming.
// A nil or empty input returns nil.
func Cap(results []Result, opts CapOptions) []Result {
	if len(results) == 0 {
		return nil
	}

	out := make([]Result, len(results))
	copy(out, results)

	if opts.OrderByKey {
		sort.Slice(out, func(i, j int) bool {
			return out[i].Key < out[j].Key
		})
	}

	if opts.MaxResults > 0 && len(out) > opts.MaxResults {
		out = out[:opts.MaxResults]
	}

	return out
}
