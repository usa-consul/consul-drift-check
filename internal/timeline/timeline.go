// Package timeline records successive drift-check results and exposes a
// time-ordered view so callers can detect trends over multiple runs.
package timeline

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/your-org/consul-drift-check/internal/diff"
)

// Entry is a single point-in-time snapshot of drift results.
type Entry struct {
	RecordedAt time.Time     `json:"recorded_at"`
	Results    []diff.Result `json:"results"`
}

// Timeline is an ordered slice of entries.
type Timeline []Entry

// Record appends a new entry to the file at path, creating it when absent.
func Record(path string, results []diff.Result) error {
	entries, _ := Load(path) // ignore missing-file error
	entries = append(entries, Entry{
		RecordedAt: time.Now().UTC(),
		Results:    results,
	})
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].RecordedAt.Before(entries[j].RecordedAt)
	})
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("timeline: create %s: %w", path, err)
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(entries)
}

// Load reads all entries from the file at path.
func Load(path string) (Timeline, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("timeline: open %s: %w", path, err)
	}
	defer f.Close()
	var t Timeline
	if err := json.NewDecoder(f).Decode(&t); err != nil {
		return nil, fmt.Errorf("timeline: decode: %w", err)
	}
	return t, nil
}

// Since returns entries recorded on or after the given time.
func (t Timeline) Since(ts time.Time) Timeline {
	var out Timeline
	for _, e := range t {
		if !e.RecordedAt.Before(ts) {
			out = append(out, e)
		}
	}
	return out
}
