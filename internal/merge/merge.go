// Package merge combines KV pair slices from multiple sources,
// applying a configurable conflict resolution strategy.
package merge

import (
	"sort"

	"github.com/hashicorp/consul/api"
)

// Strategy controls how conflicting keys are resolved.
type Strategy string

const (
	// StrategySourceWins keeps the value from the first (source) slice.
	StrategySourceWins Strategy = "source"
	// StrategyDestWins keeps the value from the second (destination) slice.
	StrategyDestWins Strategy = "destination"
	// StrategyLatest keeps the value with the higher ModifyIndex.
	StrategyLatest Strategy = "latest"
)

// Options configures the merge behaviour.
type Options struct {
	Strategy Strategy
}

// Apply merges src and dst KV pairs into a single deduplicated slice.
// Keys present in both slices are resolved according to opts.Strategy.
// The returned slice is sorted by key.
func Apply(src, dst []*api.KVPair, opts Options) []*api.KVPair {
	if opts.Strategy == "" {
		opts.Strategy = StrategySourceWins
	}

	index := make(map[string]*api.KVPair, len(src)+len(dst))

	// Load destination first so source can overwrite when StrategySourceWins.
	for _, p := range dst {
		if p == nil {
			continue
		}
		index[p.Key] = p
	}

	for _, p := range src {
		if p == nil {
			continue
		}
		existing, conflict := index[p.Key]
		if !conflict {
			index[p.Key] = p
			continue
		}
		switch opts.Strategy {
		case StrategySourceWins:
			index[p.Key] = p
		case StrategyDestWins:
			// existing (dst) already stored — nothing to do.
			_ = existing
		case StrategyLatest:
			if p.ModifyIndex > existing.ModifyIndex {
				index[p.Key] = p
			}
		}
	}

	result := make([]*api.KVPair, 0, len(index))
	for _, p := range index {
		result = append(result, p)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Key < result[j].Key
	})
	return result
}
