package diff

import (
	"sort"

	"github.com/your-org/consul-drift-check/internal/consul"
)

// TopKOptions controls how the top-K results are selected.
type TopKOptions struct {
	// K is the maximum number of results to return. Defaults to 10.
	K int
	// OrderBy controls the ranking field: "weight" (default) or "key".
	OrderBy string
}

// TopKResult holds a single ranked drift result.
type TopKResult struct {
	Rank   int
	Key    string
	Status string
	Weight float64
	Source []byte
	Dest   []byte
}

// TopK returns the K highest-weighted diff results from a weighted result set.
// Results must have been processed by diff.Apply (weight) beforehand; the
// Weight field is read directly from WeightedResult.
func TopK(results []WeightedResult, opts TopKOptions) []TopKResult {
	if opts.K <= 0 {
		opts.K = 10
	}

	sorted := make([]WeightedResult, len(results))
	copy(sorted, results)

	switch opts.OrderBy {
	case "key":
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Key < sorted[j].Key
		})
	default:
		sort.Slice(sorted, func(i, j int) bool {
			if sorted[i].Weight != sorted[j].Weight {
				return sorted[i].Weight > sorted[j].Weight
			}
			return sorted[i].Key < sorted[j].Key
		})
	}

	if opts.K < len(sorted) {
		sorted = sorted[:opts.K]
	}

	out := make([]TopKResult, len(sorted))
	for i, r := range sorted {
		out[i] = TopKResult{
			Rank:   i + 1,
			Key:    r.Key,
			Status: r.Status,
			Weight: r.Weight,
			Source: consul.CloneBytes(r.Source),
			Dest:   consul.CloneBytes(r.Dest),
		}
	}
	return out
}
