package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// HistoryEntry records a single drift-check run.
type HistoryEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Prefix    string    `json:"prefix"`
	Results   []Result  `json:"results"`
}

// AppendHistory appends entry to a newline-delimited JSON file at path.
func AppendHistory(path string, entry HistoryEntry) error {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("history: open %s: %w", path, err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	if err := enc.Encode(entry); err != nil {
		return fmt.Errorf("history: encode: %w", err)
	}
	return nil
}

// LoadHistory reads all entries from path. Returns nil if the file does not exist.
func LoadHistory(path string) ([]HistoryEntry, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("history: open %s: %w", path, err)
	}
	defer f.Close()
	var entries []HistoryEntry
	dec := json.NewDecoder(f)
	for dec.More() {
		var e HistoryEntry
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("history: decode: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// Since returns entries recorded at or after t.
func Since(entries []HistoryEntry, t time.Time) []HistoryEntry {
	var out []HistoryEntry
	for _, e := range entries {
		if !e.Timestamp.Before(t) {
			out = append(out, e)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Timestamp.Before(out[j].Timestamp)
	})
	return out
}
