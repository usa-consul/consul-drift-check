// Package window provides a sliding-window aggregator that tracks drift
// results over a fixed number of recent check cycles.
package window

import (
	"sync"

	"github.com/your-org/consul-drift-check/internal/diff"
)

// Window holds the last N snapshots of diff results.
type Window struct {
	mu      sync.Mutex
	size    int
	slots   [][]diff.Result
}

// New returns a Window that retains at most size snapshots.
// If size is less than 1 it defaults to 1.
func New(size int) *Window {
	if size < 1 {
		size = 1
	}
	return &Window{size: size}
}

// Push appends a new set of results, evicting the oldest entry when the
// window is full.
func (w *Window) Push(results []diff.Result) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if len(w.slots) >= w.size {
		w.slots = w.slots[1:]
	}
	copy_ := make([]diff.Result, len(results))
	copy(copy_, results)
	w.slots = append(w.slots, copy_)
}

// Slots returns a copy of all retained snapshots, oldest first.
func (w *Window) Slots() [][]diff.Result {
	w.mu.Lock()
	defer w.mu.Unlock()
	out := make([][]diff.Result, len(w.slots))
	for i, s := range w.slots {
		tmp := make([]diff.Result, len(s))
		copy(tmp, s)
		out[i] = tmp
	}
	return out
}

// Len returns the number of snapshots currently held.
func (w *Window) Len() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return len(w.slots)
}

// PersistentKeys returns keys that appear as drifted in every retained slot.
func (w *Window) PersistentKeys() []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	if len(w.slots) == 0 {
		return nil
	}
	counts := make(map[string]int)
	for _, slot := range w.slots {
		seen := make(map[string]struct{})
		for _, r := range slot {
			if _, ok := seen[r.Key]; !ok {
				seen[r.Key] = struct{}{}
				counts[r.Key]++
			}
		}
	}
	var keys []string
	for k, c := range counts {
		if c == len(w.slots) {
			keys = append(keys, k)
		}
	}
	return keys
}
