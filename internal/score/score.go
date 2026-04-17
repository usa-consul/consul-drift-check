// Package score computes a numeric drift severity score from diff results.
package score

import (
	"github.com/your-org/consul-drift-check/internal/diff"
)

// Weights assigned to each result kind.
const (
	WeightModified = 2.0
	WeightOnlyInSource = 1.0
	WeightOnlyInDest = 1.0
)

// Summary holds the computed score and contributing counts.
type Summary struct {
	Modified       int     `json:"modified"`
	OnlyInSource   int     `json:"only_in_source"`
	OnlyInDest     int     `json:"only_in_dest"`
	Score          float64 `json:"score"`
}

// Compute calculates a weighted drift score from a slice of diff results.
func Compute(results []diff.Result) Summary {
	var s Summary
	for _, r := range results {
		switch r.Status {
		case diff.StatusModified:
			s.Modified++
		case diff.StatusOnlyInSource:
			s.OnlyInSource++
		case diff.StatusOnlyInDestination:
			s.OnlyInDest++
		}
	}
	s.Score = float64(s.Modified)*WeightModified +
		float64(s.OnlyInSource)*WeightOnlyInSource +
		float64(s.OnlyInDest)*WeightOnlyInDest
	return s
}
