package diff

import (
	"sort"
	"strings"
)

// ScopeOptions controls how Scope filters results.
type ScopeOptions struct {
	// Prefixes limits results to keys that start with at least one of the
	// listed prefixes. An empty slice means no prefix restriction.
	Prefixes []string

	// Statuses limits results to entries whose Status matches one of the
	// listed values (case-insensitive). An empty slice means all statuses.
	Statuses []string

	// MaxResults caps the total number of returned results. Zero means
	// unlimited.
	MaxResults int
}

// Scope returns the subset of results that fall within the defined scope.
// Results are sorted by key before any MaxResults cap is applied.
func Scope(results []Result, opts ScopeOptions) []Result {
	if len(results) == 0 {
		return nil
	}

	statusSet := buildStatusSet(opts.Statuses)

	var out []Result
	for _, r := range results {
		if !matchesScopePrefixes(r.Key, opts.Prefixes) {
			continue
		}
		if len(statusSet) > 0 {
			if _, ok := statusSet[strings.ToLower(r.Status)]; !ok {
				continue
			}
		}
		out = append(out, r)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Key < out[j].Key
	})

	if opts.MaxResults > 0 && len(out) > opts.MaxResults {
		out = out[:opts.MaxResults]
	}

	return out
}

func matchesScopePrefixes(key string, prefixes []string) bool {
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

func buildStatusSet(statuses []string) map[string]struct{} {
	if len(statuses) == 0 {
		return nil
	}
	m := make(map[string]struct{}, len(statuses))
	for _, s := range statuses {
		m[strings.ToLower(s)] = struct{}{}
	}
	return m
}
