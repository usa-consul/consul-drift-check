package export_test

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"strings"
	"testing"

	"github.com/your-org/consul-drift-check/internal/diff"
	"github.com/your-org/consul-drift-check/internal/export"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "app/env", Status: diff.StatusModified, SourceValue: []byte("prod"), DestinationValue: []byte("staging")},
		{Key: "app/debug", Status: diff.StatusOnlyInSource, SourceValue: []byte("true"), DestinationValue: nil},
	}
}

func TestNewWriter_InvalidFormat(t *testing.T) {
	_, err := export.NewWriter("xml", &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestNewWriter_ValidFormats(t *testing.T) {
	for _, f := range []export.Format{export.FormatCSV, export.FormatNDJSON} {
		_, err := export.NewWriter(f, &bytes.Buffer{})
		if err != nil {
			t.Fatalf("unexpected error for format %s: %v", f, err)
		}
	}
}

func TestWrite_CSV_HasHeader(t *testing.T) {
	var buf bytes.Buffer
	w, _ := export.NewWriter(export.FormatCSV, &buf)
	if err := w.Write(sampleResults()); err != nil {
		t.Fatalf("Write: %v", err)
	}
	r := csv.NewReader(&buf)
	header, err := r.Read()
	if err != nil {
		t.Fatalf("read header: %v", err)
	}
	if header[0] != "key" || header[1] != "status" {
		t.Fatalf("unexpected header: %v", header)
	}
}

func TestWrite_CSV_RowCount(t *testing.T) {
	var buf bytes.Buffer
	w, _ := export.NewWriter(export.FormatCSV, &buf)
	_ = w.Write(sampleResults())
	r := csv.NewReader(&buf)
	rows, _ := r.ReadAll()
	// header + 2 data rows
	if len(rows) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(rows))
	}
}

func TestWrite_NDJSON_LineCount(t *testing.T) {
	var buf bytes.Buffer
	w, _ := export.NewWriter(export.FormatNDJSON, &buf)
	_ = w.Write(sampleResults())
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
}

func TestWrite_NDJSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	w, _ := export.NewWriter(export.FormatNDJSON, &buf)
	_ = w.Write(sampleResults())
	for _, line := range strings.Split(strings.TrimSpace(buf.String()), "\n") {
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			t.Fatalf("invalid JSON line %q: %v", line, err)
		}
		if _, ok := m["key"]; !ok {
			t.Fatalf("missing 'key' field in %v", m)
		}
	}
}

func TestWrite_EmptyResults_NoError(t *testing.T) {
	var buf bytes.Buffer
	w, _ := export.NewWriter(export.FormatCSV, &buf)
	if err := w.Write(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
