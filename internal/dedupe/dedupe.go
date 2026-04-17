// Package dedupe removes duplicate KV pairs from a result set,
// keeping the last occurrence when keys collide across merged sources.
package dedupe

import (
	"github.com/hashicorp/consul/api"
)

// Options controls deduplication behaviour.
type Options struct {
	// PreferSource keeps the source value when a key exists in both source and
	// destination lists. Defaults to false (last-write wins).
	PreferSource bool
}

// Apply removes duplicate keys from pairs, returning a slice where every key
// is unique. When PreferSource is false the last occurrence wins; when true
// the first occurrence wins.
func Apply(pairs []*api.KVPair, opts Options) []*api.KVPair {
	if len(pairs) == 0 {
		return pairs
	}

	seen := make(map[string]int, len(pairs)) // key -> index in out
	out := make([]*api.KVPair, 0, len(pairs))

	for _, p := range pairs {
		if idx, exists := seen[p.Key]; exists {
			if !opts.PreferSource {
				// last-write wins: replace the earlier entry
				out[idx] = p
			}
			// PreferSource: keep first, skip this entry
			continue
		}
		seen[p.Key] = len(out)
		out = append(out, p)
	}

	return out
}
