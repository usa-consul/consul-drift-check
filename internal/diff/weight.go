package diff

// Weight assigns a numeric importance score to a drift result based on its
// status and optional per-prefix multipliers.

// WeightOptions configures how scores are calculated.
type WeightOptions struct {
	// PrefixWeights maps a key prefix to a multiplier (e.g. "service/" -> 2.0).
	PrefixWeights map[string]float64
	// Base scores per status.
	ModifiedScore float64
	OnlyInSrcScore float64
	OnlyInDstScore float64
}

// DefaultWeightOptions returns sensible defaults.
func DefaultWeightOptions() WeightOptions {
	return WeightOptions{
		ModifiedScore:  2.0,
		OnlyInSrcScore: 1.0,
		OnlyInDstScore: 1.0,
		PrefixWeights:  map[string]float64{},
	}
}

// WeightedResult pairs a Result with its computed weight.
type WeightedResult struct {
	Result
	Weight float64
}

// Apply computes a weight for each result and returns a slice of WeightedResult
// sorted descending by weight.
func Apply(results []Result, opts WeightOptions) []WeightedResult {
	out := make([]WeightedResult, 0, len(results))
	for _, r := range results {
		w := baseScore(r, opts) * multiplier(r.Key, opts.PrefixWeights)
		out = append(out, WeightedResult{Result: r, Weight: w})
	}
	sortByWeight(out)
	return out
}

func baseScore(r Result, opts WeightOptions) float64 {
	switch r.Status {
	case StatusModified:
		return opts.ModifiedScore
	case StatusOnlyInSource:
		return opts.OnlyInSrcScore
	case StatusOnlyInDestination:
		return opts.OnlyInDstScore
	default:
		return 0
	}
}

func multiplier(key string, prefixWeights map[string]float64) float64 {
	best := 1.0
	for prefix, m := range prefixWeights {
		if len(prefix) > 0 && len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			if m > best {
				best = m
			}
		}
	}
	return best
}

func sortByWeight(results []WeightedResult) {
	for i := 1; i < len(results); i++ {
		for j := i; j > 0 && results[j].Weight > results[j-1].Weight; j-- {
			results[j], results[j-1] = results[j-1], results[j]
		}
	}
}
