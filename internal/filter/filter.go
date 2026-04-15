package filter

import (
	"strings"

	"github.com/hashicorp/consul/api"
)

// Options holds configuration for filtering KV entries.
type Options struct {
	// Prefixes is a list of key prefixes to include. If empty, all keys are included.
	Prefixes []string
	// ExcludeKeys is a list of exact keys to exclude from comparison.
	ExcludeKeys []string
	// ExcludePrefixes is a list of prefixes whose matching keys will be excluded.
	ExcludePrefixes []string
}

// Apply filters a slice of KV pairs according to the provided Options.
// It returns only the entries that pass all filter criteria.
func Apply(entries api.KVPairs, opts Options) api.KVPairs {
	result := make(api.KVPairs, 0, len(entries))

	for _, entry := range entries {
		if entry == nil {
			continue
		}

		if !matchesIncludedPrefixes(entry.Key, opts.Prefixes) {
			continue
		}

		if isExcludedKey(entry.Key, opts.ExcludeKeys) {
			continue
		}

		if matchesExcludedPrefix(entry.Key, opts.ExcludePrefixes) {
			continue
		}

		result = append(result, entry)
	}

	return result
}

// matchesIncludedPrefixes returns true if the key matches at least one of the
// provided prefixes, or if no prefixes are specified (include all).
func matchesIncludedPrefixes(key string, prefixes []string) bool {
	if len(prefixes) == 0 {
		return true
	}
	for _, p := range prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}

// isExcludedKey returns true if the key exactly matches one of the excluded keys.
func isExcludedKey(key string, excludeKeys []string) bool {
	for _, k := range excludeKeys {
		if key == k {
			return true
		}
	}
	return false
}

// matchesExcludedPrefix returns true if the key starts with any of the excluded prefixes.
func matchesExcludedPrefix(key string, excludePrefixes []string) bool {
	for _, p := range excludePrefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}
