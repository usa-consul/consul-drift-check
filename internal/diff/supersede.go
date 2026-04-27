package diff

import (
	"sort"
	"strings"
)

// SupersedeOptions controls which results are suppressed by higher-priority
// entries that share a common key prefix.
type SupersedeOptions struct {
	// Prefixes lists the key prefixes whose results supersede any other result
	// that shares the same top-level segment but is not covered by these prefixes.
	Prefixes []string

	// CaseSensitive disables case-folding when matching prefixes.
	CaseSensitive bool
}

// Supersede removes results whose keys are shadowed by a higher-priority prefix
// entry in the same result set. When two results share a top-level key segment
// and one of them matches a prefix listed in opts.Prefixes, the non-matching
// result is dropped.
//
// Results are returned sorted by key. If opts.Prefixes is empty the input
// slice is returned unchanged.
func Supersede(results []Result, opts SupersedeOptions) []Result {
	if len(results) == 0 || len(opts.Prefixes) == 0 {
		return results
	}

	normalise := func(s string) string {
		if opts.CaseSensitive {
			return s
		}
		return strings.ToLower(s)
	}

	// Build a set of normalised priority prefixes.
	priority := make(map[string]struct{}, len(opts.Prefixes))
	for _, p := range opts.Prefixes {
		priority[normalise(strings.TrimRight(p, "/"))] = struct{}{}
	}

	// Determine which top-level segments have at least one priority match.
	activeSeg := make(map[string]bool)
	for _, r := range results {
		nk := normalise(strings.TrimLeft(r.Key, "/"))
		for p := range priority {
			if strings.HasPrefix(nk, p+"/") || nk == p {
				seg := topSupersedeSegment(nk)
				activeSeg[seg] = true
			}
		}
	}

	out := results[:0:0]
	for _, r := range results {
		nk := normalise(strings.TrimLeft(r.Key, "/"))
		seg := topSupersedeSegment(nk)

		if !activeSeg[seg] {
			// Segment has no priority match — keep unconditionally.
			out = append(out, r)
			continue
		}

		// Only keep if this key itself matches a priority prefix.
		for p := range priority {
			if strings.HasPrefix(nk, p+"/") || nk == p {
				out = append(out, r)
				break
			}
		}
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out
}

func topSupersedeSegment(key string) string {
	if idx := strings.Index(key, "/"); idx >= 0 {
		return key[:idx]
	}
	return key
}
