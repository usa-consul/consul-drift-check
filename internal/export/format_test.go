package export_test

import (
	"testing"

	"github.com/your-org/consul-drift-check/internal/export"
)

func TestParseFormat_Valid(t *testing.T) {
	cases := []struct {
		input    string
		want     export.Format
	}{
		{"csv", export.FormatCSV},
		{"ndjson", export.FormatNDJSON},
	}
	for _, tc := range cases {
		got, err := export.ParseFormat(tc.input)
		if err != nil {
			t.Fatalf("ParseFormat(%q): unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Fatalf("ParseFormat(%q) = %q; want %q", tc.input, got, tc.want)
		}
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := export.ParseFormat("xml")
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestFormats_ReturnsAll(t *testing.T) {
	formats := export.Formats()
	if len(formats) != 2 {
		t.Fatalf("expected 2 formats, got %d", len(formats))
	}
}
