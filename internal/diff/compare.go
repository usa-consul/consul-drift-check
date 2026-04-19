package diff

import (
	"sort"
	"strings"
)

// CompareOptions controls how a prefix comparison is performed.
type CompareOptions struct {
	// CaseSensitive controls whether key comparison is case-sensitive.
	CaseSensitive bool
	// IgnorePrefixes skips keys that start with any of the given prefixes.
	IgnorePrefixes []string
}

// PrefixResult holds the comparison outcome for a single key prefix segment.
type PrefixResult struct {
	Prefix  string
	Added   int
	Removed int
	Modified int
}

// CompareByPrefix groups diff results by their top-level key prefix and
// returns a slice of PrefixResult sorted alphabetically by prefix.
func CompareByPrefix(results []Result, opts CompareOptions) []PrefixResult {
	counts := make(map[string]*PrefixResult)

	for i := range results {
		r := &results[i]
		key := r.Key
		if !opts.CaseSensitive {
			key = strings.ToLower(key)
		}
		if isIgnored(key, opts.IgnorePrefixes) {
			continue
		}
		seg := topCompareSegment(key)
		if _, ok := counts[seg]; !ok {
			counts[seg] = &PrefixResult{Prefix: seg}
		}
		switch r.Status {
		case "only_in_source":
			counts[seg].Added++
		case "only_in_destination":
			counts[seg].Removed++
		case "modified":
			counts[seg].Modified++
		}
	}

	out := make([]PrefixResult, 0, len(counts))
	for _, v := range counts {
		out = append(out, *v)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Prefix < out[j].Prefix
	})
	return out
}

func isIgnored(key string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}

func topCompareSegment(key string) string {
	key = strings.TrimPrefix(key, "/")
	if idx := strings.Index(key, "/"); idx >= 0 {
		return key[:idx]
	}
	return key
}
