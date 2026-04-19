package diff

import (
	"math/rand"
	"sort"
)

// SampleOptions controls how results are sampled.
type SampleOptions struct {
	// MaxResults is the maximum number of results to return.
	// If zero or negative, all results are returned.
	MaxResults int
	// Seed is used to initialise the random source for reproducibility.
	// If zero, sampling is non-deterministic.
	Seed int64
	// Deterministic sorts by key before sampling so that, for a given
	// seed, the same subset is always chosen.
	Deterministic bool
}

// Sample returns a random subset of results according to opts.
// The returned slice is sorted by key.
func Sample(results []Result, opts SampleOptions) []Result {
	if len(results) == 0 {
		return nil
	}
	if opts.MaxResults <= 0 || opts.MaxResults >= len(results) {
		out := make([]Result, len(results))
		copy(out, results)
		sortByKey(out)
		return out
	}

	pool := make([]Result, len(results))
	copy(pool, results)

	if opts.Deterministic {
		sortByKey(pool)
	}

	var r *rand.Rand
	if opts.Seed != 0 {
		//nolint:gosec
		r = rand.New(rand.NewSource(opts.Seed))
	} else {
		//nolint:gosec
		r = rand.New(rand.NewSource(rand.Int63()))
	}

	r.Shuffle(len(pool), func(i, j int) {
		pool[i], pool[j] = pool[j], pool[i]
	})

	sampled := pool[:opts.MaxResults]
	sort.Slice(sampled, func(i, j int) bool {
		return sampled[i].Key < sampled[j].Key
	})
	return sampled
}
