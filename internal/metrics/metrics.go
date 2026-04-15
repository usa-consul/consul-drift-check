package metrics

import (
	"time"

	"github.com/your-org/consul-drift-check/internal/diff"
)

// Summary holds aggregated statistics for a drift check run.
type Summary struct {
	TotalKeys       int           `json:"total_keys"`
	MatchedKeys     int           `json:"matched_keys"`
	OnlyInSource    int           `json:"only_in_source"`
	OnlyInDest      int           `json:"only_in_destination"`
	ModifiedKeys    int           `json:"modified_keys"`
	DriftDetected   bool          `json:"drift_detected"`
	Duration        time.Duration `json:"duration_ms"`
	CheckedAt       time.Time     `json:"checked_at"`
}

// Collect builds a Summary from a slice of diff results and the elapsed duration.
func Collect(results []diff.Result, duration time.Duration) Summary {
	s := Summary{
		Duration:  duration,
		CheckedAt: time.Now().UTC(),
	}

	for _, r := range results {
		switch r.Status {
		case diff.StatusMatch:
			s.MatchedKeys++
		case diff.StatusOnlyInSource:
			s.OnlyInSource++
		case diff.StatusOnlyInDest:
			s.OnlyInDest++
		case diff.StatusModified:
			s.ModifiedKeys++
		}
	}

	s.TotalKeys = s.MatchedKeys + s.OnlyInSource + s.OnlyInDest + s.ModifiedKeys
	s.DriftDetected = s.OnlyInSource > 0 || s.OnlyInDest > 0 || s.ModifiedKeys > 0

	return s
}
