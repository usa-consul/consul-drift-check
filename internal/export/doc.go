// Package export provides writers that serialise drift results into
// portable file formats for downstream processing.
//
// Supported formats:
//
//	"csv"    — comma-separated values with a header row.
//	"ndjson" — newline-delimited JSON, one object per line.
//
// Usage:
//
//	w, err := export.NewWriter(export.FormatCSV, os.Stdout)
//	if err != nil { ... }
//	err = w.Write(results)
package export
