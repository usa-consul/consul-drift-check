package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/your-org/consul-drift-check/internal/diff"
)

// Snapshot represents a point-in-time capture of drift results.
type Snapshot struct {
	Timestamp time.Time         `json:"timestamp"`
	Prefix    string            `json:"prefix"`
	Results   []diff.DiffResult `json:"results"`
}

// Save writes a snapshot to the given file path as JSON.
func Save(path string, prefix string, results []diff.DiffResult) error {
	s := Snapshot{
		Timestamp: time.Now().UTC(),
		Prefix:    prefix,
		Results:   results,
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal failed: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write failed: %w", err)
	}

	return nil
}

// Load reads a snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read failed: %w", err)
	}

	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal failed: %w", err)
	}

	return &s, nil
}
