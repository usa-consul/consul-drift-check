// Package groupby groups diff results by a key segment or label.
package groupby

import (
	"strings"

	"github.com/your-org/consul-drift-check/internal/diff"
)

// Options controls how results are grouped.
type Options struct {
	// SegmentIndex is the zero-based path segment used as the group key.
	// e.g. for "services/web/port", index 0 yields "services".
	SegmentIndex int
	// Separator splits the KV key into segments. Defaults to "/".
	Separator string
}

// Group maps a group key to the slice of results that belong to it.
type Group map[string][]diff.Result

// Apply partitions results into groups based on a path segment.
func Apply(results []diff.Result, opts Options) Group {
	sep := opts.Separator
	if sep == "" {
		sep = "/"
	}

	groups := make(Group)
	for _, r := range results {
		key := groupKey(r.Key, sep, opts.SegmentIndex)
		groups[key] = append(groups[key], r)
	}
	return groups
}

// Keys returns the sorted group keys present in g.
func Keys(g Group) []string {
	out := make([]string, 0, len(g))
	for k := range g {
		out = append(out, k)
	}
	sortStrings(out)
	return out
}

func groupKey(key, sep string, idx int) string {
	parts := strings.Split(strings.TrimPrefix(key, sep), sep)
	if idx < len(parts) {
		return parts[idx]
	}
	return "_other"
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
