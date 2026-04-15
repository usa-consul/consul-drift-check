package report_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/example/consul-drift-check/internal/diff"
	"github.com/example/consul-drift-check/internal/report"
)

func TestWriter_TextFormat_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	w := report.NewWriter(&buf, report.FormatText)

	if err := w.Write([]diff.Result{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "No drift detected.") {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

func TestWriter_TextFormat_WithResults(t *testing.T) {
	var buf bytes.Buffer
	w := report.NewWriter(&buf, report.FormatText)

	results := []diff.Result{
		{Key: "config/app/port", Status: diff.StatusOnlyInSource},
		{Key: "config/app/host", Status: diff.StatusModified},
	}

	if err := w.Write(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "config/app/port") {
		t.Errorf("expected key in output, got: %s", out)
	}
	if !strings.Contains(out, "config/app/host") {
		t.Errorf("expected key in output, got: %s", out)
	}
}

func TestWriter_JSONFormat_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	w := report.NewWriter(&buf, report.FormatJSON)

	if err := w.Write([]diff.Result{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.TrimSpace(buf.String()) != "[]" {
		t.Errorf("expected empty JSON array, got: %s", buf.String())
	}
}

func TestWriter_JSONFormat_WithResults(t *testing.T) {
	var buf bytes.Buffer
	w := report.NewWriter(&buf, report.FormatJSON)

	results := []diff.Result{
		{Key: "config/db/url", Status: diff.StatusOnlyInDestination},
	}

	if err := w.Write(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "config/db/url") {
		t.Errorf("expected key in JSON output, got: %s", out)
	}
	if !strings.Contains(out, "only_in_destination") {
		t.Errorf("expected status in JSON output, got: %s", out)
	}
}
