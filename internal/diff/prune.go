package diff

import (
	"strings"

	"github.com/hashicorp/consul/api"
)

// PruneOptions controls which results are dropped by Prune.
type PruneOptions struct {
	// MaxAgeDays drops results whose key has not changed across snapshots
	// when the key matches any of these prefixes. Zero means no age filter.
	StalePrefixes []string

	// EmptyValues drops results where both source and destination values are empty.
	DropEmptyValues bool

	// OnlyStatus retains only results whose status matches one of the given
	// values (case-insensitive). Empty slice means retain all.
	OnlyStatus []string
}

// Prune removes unwanted entries from results according to opts.
func Prune(results []*api.KVPair, src, dst map[string]*api.KVPair, opts PruneOptions) []*api.KVPair {
	if len(results) == 0 {
		return nil
	}

	statusSet := make(map[string]struct{}, len(opts.OnlyStatus))
	for _, s := range opts.OnlyStatus {
		statusSet[strings.ToLower(s)] = struct{}{}
	}

	out := make([]*api.KVPair, 0, len(results))
	for _, r := range results {
		if r == nil {
			continue
		}

		if opts.DropEmptyValues {
			sv := src[r.Key]
			dv := dst[r.Key]
			sEmpty := sv == nil || len(sv.Value) == 0
			dEmpty := dv == nil || len(dv.Value) == 0
			if sEmpty && dEmpty {
				continue
			}
		}

		if len(statusSet) > 0 {
			// Status is stored in the Flags field by convention in this project:
			// we derive status label from presence in src/dst maps.
			status := deriveStatus(r.Key, src, dst)
			if _, ok := statusSet[status]; !ok {
				continue
			}
		}

		if matchesStalePrefixes(r.Key, opts.StalePrefixes) {
			continue
		}

		out = append(out, r)
	}

	if len(out) == 0 {
		return nil
	}
	return out
}

func deriveStatus(key string, src, dst map[string]*api.KVPair) string {
	_, inSrc := src[key]
	_, inDst := dst[key]
	switch {
	case inSrc && !inDst:
		return "source_only"
	case !inSrc && inDst:
		return "dest_only"
	default:
		return "modified"
	}
}

func matchesStalePrefixes(key string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}
