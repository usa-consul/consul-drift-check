package diff

import (
	"bytes"
	"sort"

	"github.com/your-org/consul-drift-check/internal/consul"
)

// Status represents the type of drift detected for a key.
type Status string

const (
	OnlyInSource      Status = "only_in_source"
	OnlyInDestination Status = "only_in_destination"
	Modified          Status = "modified"
	InSync            Status = "in_sync"
)

// DiffResult holds the comparison outcome for a single KV key.
type DiffResult struct {
	Key         string  `json:"key"`
	Status      Status  `json:"status"`
	SourceValue []byte  `json:"source_value,omitempty"`
	DestValue   []byte  `json:"dest_value,omitempty"`
}

// Compare computes the diff between source and destination KV entry maps.
func Compare(source, dest map[string]*consul.KVEntry) []DiffResult {
	var results []DiffResult

	for key, srcEntry := range source {
		dstEntry, exists := dest[key]
		if !exists {
			results = append(results, DiffResult{
				Key:         key,
				Status:      OnlyInSource,
				SourceValue: srcEntry.Value,
			})
			continue
		}
		if !bytesEqual(srcEntry.Value, dstEntry.Value) {
			results = append(results, DiffResult{
				Key:         key,
				Status:      Modified,
				SourceValue: srcEntry.Value,
				DestValue:   dstEntry.Value,
			})
		}
	}

	for key, dstEntry := range dest {
		if _, exists := source[key]; !exists {
			results = append(results, DiffResult{
				Key:       key,
				Status:    OnlyInDestination,
				DestValue: dstEntry.Value,
			})
		}
	}

	sortByKey(results)
	return results
}

func bytesEqual(a, b []byte) bool {
	return bytes.Equal(a, b)
}

func sortByKey(results []DiffResult) {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})
}
