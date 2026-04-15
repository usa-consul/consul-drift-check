// Package alert provides threshold-based alerting for drift check results.
package alert

import (
	"fmt"
	"io"

	"github.com/example/consul-drift-check/internal/metrics"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelOK      Level = "OK"
	LevelWarning Level = "WARNING"
	LevelCritical Level = "CRITICAL"
)

// Thresholds defines the drift count limits for each alert level.
type Thresholds struct {
	Warning  int
	Critical int
}

// Result holds the evaluated alert level and a human-readable message.
type Result struct {
	Level   Level
	Message string
}

// Evaluate compares collected metrics against thresholds and returns an alert Result.
func Evaluate(m metrics.Summary, t Thresholds) Result {
	total := m.Added + m.Removed + m.Modified

	switch {
	case total >= t.Critical:
		return Result{
			Level:   LevelCritical,
			Message: fmt.Sprintf("drift count %d exceeds critical threshold %d", total, t.Critical),
		}
	case total >= t.Warning:
		return Result{
			Level:   LevelWarning,
			Message: fmt.Sprintf("drift count %d exceeds warning threshold %d", total, t.Warning),
		}
	default:
		return Result{
			Level:   LevelOK,
			Message: fmt.Sprintf("drift count %d is within acceptable limits", total),
		}
	}
}

// Write prints the alert result to w in a human-readable format.
func Write(w io.Writer, r Result) {
	fmt.Fprintf(w, "[%s] %s\n", r.Level, r.Message)
}
