package report

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/example/consul-drift-check/internal/diff"
)

// Format represents the output format for drift reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Writer writes drift comparison results to an output stream.
type Writer struct {
	out    io.Writer
	format Format
}

// NewWriter creates a new report Writer with the given output and format.
func NewWriter(out io.Writer, format Format) *Writer {
	return &Writer{out: out, format: format}
}

// Write outputs the diff results according to the configured format.
func (w *Writer) Write(results []diff.Result) error {
	switch w.format {
	case FormatJSON:
		return w.writeJSON(results)
	default:
		return w.writeText(results)
	}
}

func (w *Writer) writeText(results []diff.Result) error {
	if len(results) == 0 {
		fmt.Fprintln(w.out, "No drift detected.")
		return nil
	}

	tw := tabwriter.NewWriter(w.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "STATUS\tKEY")
	for _, r := range results {
		fmt.Fprintf(tw, "%s\t%s\n", r.Status, r.Key)
	}
	return tw.Flush()
}

func (w *Writer) writeJSON(results []diff.Result) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w.out, "[]")
		return err
	}

	fmt.Fprintln(w.out, "[")
	for i, r := range results {
		comma := ","
		if i == len(results)-1 {
			comma = ""
		}
		fmt.Fprintf(w.out, "  {\"status\": \"%s\", \"key\": \"%s\"}%s\n", r.Status, r.Key, comma)
	}
	_, err := fmt.Fprintln(w.out, "]")
	return err
}
