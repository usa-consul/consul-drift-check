// Package metrics collects aggregate statistics from drift check results.
package metrics

import "github.com/example/consul-drift-check/internal/diff"

// Summary holds aggregate counts derived from a set of diff results.
type Summary struct {
	Added    int
	Removed  int
	Modified int
	Total    int
}

// Collect aggregates diff.Result entries into a Summary.
func Collect(results []diff.Result) Summary {
	var s Summary
	for _, r := range results {
		switch r.Status {
		case diff.StatusAdded:
			s.Added++
		case diff.StatusRemoved:
			s.Removed++
		case diff.StatusModified:
			s.Modified++
		}
	}
	s.Total = s.Added + s.Removed + s.Modified
	return s
}
