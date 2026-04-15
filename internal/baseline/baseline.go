package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/your-org/consul-drift-check/internal/diff"
)

// Baseline represents a saved reference state of KV entries.
type Baseline struct {
	CreatedAt time.Time          `json:"created_at"`
	Prefix    string             `json:"prefix"`
	Entries   map[string][]byte  `json:"entries"`
}

// Save writes the given KV entries as a baseline file to path.
func Save(path, prefix string, entries map[string][]byte) error {
	b := Baseline{
		CreatedAt: time.Now().UTC(),
		Prefix:    prefix,
		Entries:   entries,
	}
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("baseline: write file: %w", err)
	}
	return nil
}

// Load reads a baseline file from path and returns the Baseline.
func Load(path string) (*Baseline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("baseline: read file: %w", err)
	}
	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("baseline: unmarshal: %w", err)
	}
	return &b, nil
}

// Compare diffs the baseline entries against current entries using the diff package.
func Compare(b *Baseline, current map[string][]byte) []diff.Result {
	return diff.Compare(b.Entries, current)
}
