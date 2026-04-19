package diff

import "sort"

// DeltaEntry represents the change in drift count for a single key prefix
// between two consecutive snapshots.
type DeltaEntry struct {
	Prefix string
	Before int
	After  int
	Delta  int // positive = more drift, negative = less drift
}

// ComputeDelta compares two rollup maps (prefix -> count) and returns a
// sorted slice of DeltaEntry describing how drift has changed per prefix.
func ComputeDelta(before, after map[string]int) []DeltaEntry {
	keys := make(map[string]struct{})
	for k := range before {
		keys[k] = struct{}{}
	}
	for k := range after {
		keys[k] = struct{}{}
	}

	entries := make([]DeltaEntry, 0, len(keys))
	for k := range keys {
		b := before[k]
		a := after[k]
		entries = append(entries, DeltaEntry{
			Prefix: k,
			Before: b,
			After:  a,
			Delta:  a - b,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Delta != entries[j].Delta {
			return entries[i].Delta > entries[j].Delta
		}
		return entries[i].Prefix < entries[j].Prefix
	})

	return entries
}

// TotalDelta returns the net change in total drift count between two rollups.
func TotalDelta(before, after map[string]int) int {
	var b, a int
	for _, v := range before {
		b += v
	}
	for _, v := range after {
		a += v
	}
	return a - b
}
