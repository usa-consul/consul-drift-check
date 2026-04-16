// Package summarize produces a human-readable drift summary from diff results.
package summarize

import (
	"fmt"
	"strings"

	"github.com/your-org/consul-drift-check/internal/diff"
)

// Summary holds aggregated statistics for a drift check run.
type Summary struct {
	Total    int
	Added    int
	Removed  int
	Modified int
	Lines    []string
}

// Build creates a Summary from a slice of diff results.
func Build(results []diff.Result) Summary {
	s := Summary{}
	for _, r := range results {
		s.Total++
		switch r.Status {
		case diff.StatusOnlyInSource:
			s.Added++
			s.Lines = append(s.Lines, fmt.Sprintf("+ %s", r.Key))
		case diff.StatusOnlyInDestination:
			s.Removed++
			s.Lines = append(s.Lines, fmt.Sprintf("- %s", r.Key))
		case diff.StatusModified:
			s.Modified++
			s.Lines = append(s.Lines, fmt.Sprintf("~ %s", r.Key))
		}
	}
	return s
}

// String returns a compact one-line representation of the summary.
func (s Summary) String() string {
	return fmt.Sprintf("total=%d added=%d removed=%d modified=%d", s.Total, s.Added, s.Removed, s.Modified)
}

// Report returns a multi-line report of all changed keys.
func (s Summary) Report() string {
	if len(s.Lines) == 0 {
		return "no drift detected"
	}
	return strings.Join(s.Lines, "\n")
}
