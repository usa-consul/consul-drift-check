package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/example/consul-drift-check/internal/alert"
	"github.com/example/consul-drift-check/internal/metrics"
)

func makeSummary(added, removed, modified int) metrics.Summary {
	return metrics.Summary{
		Added:    added,
		Removed:  removed,
		Modified: modified,
		Total:    added + removed + modified,
	}
}

func TestEvaluate_OK(t *testing.T) {
	r := alert.Evaluate(makeSummary(0, 0, 0), alert.Thresholds{Warning: 5, Critical: 10})
	if r.Level != alert.LevelOK {
		t.Fatalf("expected OK, got %s", r.Level)
	}
}

func TestEvaluate_Warning(t *testing.T) {
	r := alert.Evaluate(makeSummary(3, 2, 0), alert.Thresholds{Warning: 5, Critical: 10})
	if r.Level != alert.LevelWarning {
		t.Fatalf("expected WARNING, got %s", r.Level)
	}
}

func TestEvaluate_Critical(t *testing.T) {
	r := alert.Evaluate(makeSummary(5, 4, 2), alert.Thresholds{Warning: 5, Critical: 10})
	if r.Level != alert.LevelCritical {
		t.Fatalf("expected CRITICAL, got %s", r.Level)
	}
}

func TestEvaluate_ExactWarningBoundary(t *testing.T) {
	r := alert.Evaluate(makeSummary(5, 0, 0), alert.Thresholds{Warning: 5, Critical: 10})
	if r.Level != alert.LevelWarning {
		t.Fatalf("expected WARNING at exact boundary, got %s", r.Level)
	}
}

func TestEvaluate_MessageContainsCount(t *testing.T) {
	r := alert.Evaluate(makeSummary(2, 1, 0), alert.Thresholds{Warning: 5, Critical: 10})
	if !strings.Contains(r.Message, "3") {
		t.Errorf("expected message to contain drift count, got: %s", r.Message)
	}
}

func TestWrite_FormatsOutput(t *testing.T) {
	var buf bytes.Buffer
	alert.Write(&buf, alert.Result{Level: alert.LevelWarning, Message: "some drift detected"})
	out := buf.String()
	if !strings.Contains(out, "[WARNING]") {
		t.Errorf("expected [WARNING] in output, got: %s", out)
	}
	if !strings.Contains(out, "some drift detected") {
		t.Errorf("expected message in output, got: %s", out)
	}
}
