// Package audit records drift check runs to an append-only JSONL log.
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/example/consul-drift-check/internal/metrics"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time       `json:"timestamp"`
	ConfigPath string        `json:"config_path"`
	Prefix     string        `json:"prefix"`
	Summary    metrics.Summary `json:"summary"`
	Error      string        `json:"error,omitempty"`
}

// Logger writes audit entries to a file.
type Logger struct {
	path string
}

// NewLogger creates a Logger that appends to the given file path.
func NewLogger(path string) *Logger {
	return &Logger{path: path}
}

// Record appends an Entry to the audit log file.
func (l *Logger) Record(e Entry) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}

	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("audit: open log file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	if err := enc.Encode(e); err != nil {
		return fmt.Errorf("audit: encode entry: %w", err)
	}
	return nil
}

// ReadAll reads all entries from the audit log file.
func ReadAll(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("audit: open log file: %w", err)
	}
	defer f.Close()

	var entries []Entry
	dec := json.NewDecoder(f)
	for dec.More() {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("audit: decode entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}
