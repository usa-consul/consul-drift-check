// Package highlight marks diff results that exceed a change threshold,
// making it easy to surface high-churn namespaces in reports.
package highlight

import (
	"strings"

	"github.com/your-org/consul-drift-check/internal/diff"
)

// Level describes how prominently a result should be highlighted.
type Level string

const (
	LevelNone     Level = "none"
	LevelModified Level = "modified"
	LevelCritical Level = "critical"
)

// Result wraps a diff.Result with a highlight level.
type Result struct {
	diff.Result
	Level Level
}

// Options controls the thresholds used when assigning levels.
type Options struct {
	// CriticalPrefixes marks any key whose prefix matches as critical regardless of status.
	CriticalPrefixes []string
	// ModifiedOnly restricts highlighting to modified keys only.
	ModifiedOnly bool
}

// Apply assigns a highlight Level to each diff.Result.
func Apply(results []diff.Result, opts Options) []Result {
	out := make([]Result, 0, len(results))
	for _, r := range results {
		out = append(out, Result{
			Result: r,
			Level:  resolve(r, opts),
		})
	}
	return out
}

func resolve(r diff.Result, opts Options) Level {
	for _, p := range opts.CriticalPrefixes {
		if strings.HasPrefix(r.Key, p) {
			return LevelCritical
		}
	}
	if r.Status == diff.StatusModified {
		return LevelModified
	}
	if !opts.ModifiedOnly && (r.Status == diff.StatusOnlyInSource || r.Status == diff.StatusOnlyInDestination) {
		return LevelModified
	}
	return LevelNone
}
