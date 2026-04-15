package diff

import "sort"

// KeyValue represents a single Consul KV entry.
type KeyValue struct {
	Key   string
	Value []byte
}

// Result holds the drift comparison result between two namespaces.
type Result struct {
	OnlyInSource      []KeyValue
	OnlyInDestination []KeyValue
	Modified          []KeyValue
}

// HasDrift returns true if any differences were found.
func (r *Result) HasDrift() bool {
	return len(r.OnlyInSource) > 0 ||
		len(r.OnlyInDestination) > 0 ||
		len(r.Modified) > 0
}

// Compare detects drift between two maps of KV entries.
// source and destination are maps of key -> value.
func Compare(source, destination map[string][]byte) *Result {
	result := &Result{}

	for key, srcVal := range source {
		dstVal, exists := destination[key]
		if !exists {
			result.OnlyInSource = append(result.OnlyInSource, KeyValue{Key: key, Value: srcVal})
			continue
		}
		if !bytesEqual(srcVal, dstVal) {
			result.Modified = append(result.Modified, KeyValue{Key: key, Value: srcVal})
		}
	}

	for key, dstVal := range destination {
		if _, exists := source[key]; !exists {
			result.OnlyInDestination = append(result.OnlyInDestination, KeyValue{Key: key, Value: dstVal})
		}
	}

	sortByKey(result.OnlyInSource)
	sortByKey(result.OnlyInDestination)
	sortByKey(result.Modified)

	return result
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func sortByKey(entries []KeyValue) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
}
