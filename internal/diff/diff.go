package diff

import (
	"bytes"
	"sort"

	"github.com/your-org/consul-drift-check/internal/consul"
)

const (
	StatusOnlyInSource      = "only_in_source"
	StatusOnlyInDestination = "only_in_destination"
	StatusModified          = "modified"
	StatusMatch             = "match"
)

// Result describes the drift status of a single KV key.
type Result struct {
	Key         string `json:"key"`
	Status      string `json:"status"`
	SourceValue []byte `json:"source_value,omitempty"`
	DestValue   []byte `json:"dest_value,omitempty"`
}

// Compare returns drift results between src and dst KV maps.
func Compare(src, dst map[string]*consul.KVPair) []Result {
	keys := unionKeys(src, dst)
	var results []Result
	for _, k := range keys {
		sv, inSrc := src[k]
		dv, inDst := dst[k]
		switch {
		case inSrc && !inDst:
			results = append(results, Result{Key: k, Status: StatusOnlyInSource, SourceValue: sv.Value})
		case !inSrc && inDst:
			results = append(results, Result{Key: k, Status: StatusOnlyInDestination, DestValue: dv.Value})
		case !bytesEqual(sv.Value, dv.Value):
			results = append(results, Result{Key: k, Status: StatusModified, SourceValue: sv.Value, DestValue: dv.Value})
		default:
			results = append(results, Result{Key: k, Status: StatusMatch})
		}
	}
	return results
}

func bytesEqual(a, b []byte) bool {
	return bytes.Equal(a, b)
}

func unionKeys(src, dst map[string]*consul.KVPair) []string {
	seen := make(map[string]struct{}, len(src)+len(dst))
	for k := range src {
		seen[k] = struct{}{}
	}
	for k := range dst {
		seen[k] = struct{}{}
	}
	return sortByKey(seen)
}

func sortByKey(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
