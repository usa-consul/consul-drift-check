package diff

import (
	"bytes"
	"sort"
)

// Status represents the drift state of a single KV entry.
type Status string

const (
	StatusOnlyInSource      Status = "only_in_source"
	StatusOnlyInDestination Status = "only_in_destination"
	StatusModified          Status = "modified"
)

// Result holds the drift information for a single key.
type Result struct {
	Key    string
	Status Status
}

// Compare compares two KV maps (source and destination) and returns
// a slice of Results describing any detected drift.
func Compare(source, destination map[string][]byte) []Result {
	var results []Result

	for key, srcVal := range source {
		dstVal, exists := destination[key]
		if !exists {
			results = append(results, Result{Key: key, Status: StatusOnlyInSource})
			continue
		}
		if !bytesEqual(srcVal, dstVal) {
			results = append(results, Result{Key: key, Status: StatusModified})
		}
	}

	for key := range destination {
		if _, exists := source[key]; !exists {
			results = append(results, Result{Key: key, Status: StatusOnlyInDestination})
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
