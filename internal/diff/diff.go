// Package diff compares two sets of Consul KV entries and returns a list of
// changes between source and destination namespaces.
package diff

import (
	"bytes"
	"sort"

	"github.com/your-org/consul-drift-check/internal/consul"
)

// Status describes the kind of difference found for a key.
type Status string

const (
	StatusEqual             Status = "equal"
	StatusOnlyInSource      Status = "only_in_source"
	StatusOnlyInDestination Status = "only_in_destination"
	StatusModified          Status = "modified"
)

// Result represents a single key comparison outcome.
type Result struct {
	Key    string
	Status Status
	Source []byte
	Dest   []byte
}

// Compare returns the diff between src and dst KV maps.
func Compare(src, dst map[string]consul.KVPair) []Result {
	var results []Result

	for key, sv := range src {
		dv, ok := dst[key]
		if !ok {
			results = append(results, Result{Key: key, Status: StatusOnlyInSource, Source: sv.Value})
			continue
		}
		if !bytesEqual(sv.Value, dv.Value) {
			results = append(results, Result{Key: key, Status: StatusModified, Source: sv.Value, Dest: dv.Value})
		}
	}

	for key, dv := range dst {
		if _, ok := src[key]; !ok {
			results = append(results, Result{Key: key, Status: StatusOnlyInDestination, Dest: dv.Value})
		}
	}

	sortByKey(results)
	return results
}

func bytesEqual(a, b []byte) bool {
	return bytes.Equal(a, b)
}

func sortByKey(results []Result) {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})
}
