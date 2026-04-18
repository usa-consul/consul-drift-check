// Package sample provides periodic KV sampling for drift trending.
package sample

import (
	"encoding/json"
	"os"
	"time"

	"github.com/hashicorp/consul/api"
)

// Entry records a single KV sample taken at a point in time.
type Entry struct {
	Timestamp time.Time          `json:"timestamp"`
	Prefix    string             `json:"prefix"`
	Pairs     []*api.KVPair     `json:"pairs"`
}

// Record appends a sample entry to the file at path.
func Record(path, prefix string, pairs []*api.KVPair) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	entry := Entry{
		Timestamp: time.Now().UTC(),
		Prefix:    prefix,
		Pairs:     pairs,
	}
	return json.NewEncoder(f).Encode(entry)
}

// LoadAll reads all sample entries from path.
// Returns nil slice when the file does not exist.
func LoadAll(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var entries []Entry
	dec := json.NewDecoder(f)
	for dec.More() {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// Since returns entries recorded at or after the given time.
func Since(entries []Entry, t time.Time) []Entry {
	var out []Entry
	for _, e := range entries {
		if !e.Timestamp.Before(t) {
			out = append(out, e)
		}
	}
	return out
}
