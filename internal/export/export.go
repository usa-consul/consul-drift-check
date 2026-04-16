// Package export writes drift results to external formats (CSV, NDJSON).
package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/your-org/consul-drift-check/internal/diff"
)

// Format identifies the output format.
type Format string

const (
	FormatCSV    Format = "csv"
	FormatNDJSON Format = "ndjson"
)

// Writer exports diff results to an io.Writer.
type Writer struct {
	format Format
	out    io.Writer
}

// NewWriter returns a Writer for the given format.
func NewWriter(format Format, out io.Writer) (*Writer, error) {
	switch format {
	case FormatCSV, FormatNDJSON:
		return &Writer{format: format, out: out}, nil
	default:
		return nil, fmt.Errorf("unsupported export format: %s", format)
	}
}

// Write encodes results to the underlying writer.
func (w *Writer) Write(results []diff.Result) error {
	switch w.format {
	case FormatCSV:
		return writeCSV(w.out, results)
	case FormatNDJSON:
		return writeNDJSON(w.out, results)
	default:
		return fmt.Errorf("unknown format: %s", w.format)
	}
}

func writeCSV(out io.Writer, results []diff.Result) error {
	cw := csv.NewWriter(out)
	if err := cw.Write([]string{"key", "status", "source_value", "destination_value", "exported_at"}); err != nil {
		return err
	}
	now := time.Now().UTC().Format(time.RFC3339)
	for _, r := range results {
		if err := cw.Write([]string{r.Key, string(r.Status), string(r.SourceValue), string(r.DestinationValue), now}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

type ndjsonRow struct {
	Key              string `json:"key"`
	Status           string `json:"status"`
	SourceValue      string `json:"source_value"`
	DestinationValue string `json:"destination_value"`
	ExportedAt       string `json:"exported_at"`
}

func writeNDJSON(out io.Writer, results []diff.Result) error {
	enc := json.NewEncoder(out)
	now := time.Now().UTC().Format(time.RFC3339)
	for _, r := range results {
		row := ndjsonRow{
			Key:              r.Key,
			Status:           string(r.Status),
			SourceValue:      string(r.SourceValue),
			DestinationValue: string(r.DestinationValue),
			ExportedAt:       now,
		}
		if err := enc.Encode(row); err != nil {
			return err
		}
	}
	return nil
}
