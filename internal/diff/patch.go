package diff

import (
	"sort"

	"github.com/hashicorp/consul/api"
)

// PatchEntry describes a single remediation action to bring destination
// in line with source.
type PatchEntry struct {
	Key    string
	Action string // "set", "delete"
	Value  []byte
}

// Patch derives the ordered list of changes needed to reconcile destination
// with source based on a prior Compare result.
func Patch(results []Result) []PatchEntry {
	var entries []PatchEntry

	for _, r := range results {
		switch r.Status {
		case StatusOnlyInSource, StatusModified:
			entries = append(entries, PatchEntry{
				Key:    r.Key,
				Action: "set",
				Value:  r.SourceValue,
			})
		case StatusOnlyInDestination:
			entries = append(entries, PatchEntry{
				Key:    r.Key,
				Action: "delete",
			})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return entries
}

// ToKVPairs converts patch entries with action "set" into Consul KV pairs
// suitable for batch write operations.
func ToKVPairs(entries []PatchEntry) []*api.KVPair {
	var pairs []*api.KVPair
	for _, e := range entries {
		if e.Action == "set" {
			pairs = append(pairs, &api.KVPair{
				Key:   e.Key,
				Value: e.Value,
			})
		}
	}
	return pairs
}
