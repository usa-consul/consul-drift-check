package diff

import (
	"strings"

	"github.com/hashicorp/consul/api"
)

// ExcludeOptions controls which diff results are suppressed.
type ExcludeOptions struct {
	// Statuses lists result statuses to drop (e.g. "only_in_source").
	Statuses []string
	// Prefixes lists key prefixes whose results should be dropped.
	Prefixes []string
}

// Exclude filters out diff results that match the given options.
// It returns a new slice without modifying the original.
func Exclude(results []Result, opts ExcludeOptions) []Result {
	statusSet := make(map[string]struct{}, len(opts.Statuses))
	for _, s := range opts.Statuses {
		statusSet[strings.ToLower(s)] = struct{}{}
	}

	out := make([]Result, 0, len(results))
	for _, r := range results {
		if _, blocked := statusSet[strings.ToLower(r.Status)]; blocked {
			continue
		}
		if matchesExcludePrefix(r.Key, opts.Prefixes) {
			continue
		}
		out = append(out, r)
	}
	return out
}

// ExcludeFromPairs removes KV pairs whose keys match any of the given prefixes.
func ExcludeFromPairs(pairs []*api.KVPair, prefixes []string) []*api.KVPair {
	if len(prefixes) == 0 {
		return pairs
	}
	out := make([]*api.KVPair, 0, len(pairs))
	for _, p := range pairs {
		if p == nil || matchesExcludePrefix(p.Key, prefixes) {
			continue
		}
		out = append(out, p)
	}
	return out
}

func matchesExcludePrefix(key string, prefixes []string) bool {
	for _, pfx := range prefixes {
		if strings.HasPrefix(key, pfx) {
			return true
		}
	}
	return false
}
