// Package watch provides key-level change detection between two snapshots.
package watch

import (
	"github.com/hashicorp/consul/api"
)

// ChangeKind describes the type of change detected.
type ChangeKind string

const (
	Added    ChangeKind = "added"
	Removed  ChangeKind = "removed"
	Modified ChangeKind = "modified"
)

// Change represents a single key-level change between two snapshots.
type Change struct {
	Key      string
	Kind     ChangeKind
	OldValue []byte
	NewValue []byte
}

// Diff compares two KV pair slices and returns the list of changes.
func Diff(prev, curr api.KVPairs) []Change {
	prevMap := index(prev)
	currMap := index(curr)

	var changes []Change

	for key, oldPair := range prevMap {
		if newPair, ok := currMap[key]; !ok {
			changes = append(changes, Change{
				Key:      key,
				Kind:     Removed,
				OldValue: oldPair.Value,
			})
		} else if string(oldPair.Value) != string(newPair.Value) {
			changes = append(changes, Change{
				Key:      key,
				Kind:     Modified,
				OldValue: oldPair.Value,
				NewValue: newPair.Value,
			})
		}
	}

	for key, newPair := range currMap {
		if _, ok := prevMap[key]; !ok {
			changes = append(changes, Change{
				Key:      key,
				Kind:     Added,
				NewValue: newPair.Value,
			})
		}
	}

	return changes
}

func index(pairs api.KVPairs) map[string]*api.KVPair {
	m := make(map[string]*api.KVPair, len(pairs))
	for _, p := range pairs {
		if p != nil {
			m[p.Key] = p
		}
	}
	return m
}
