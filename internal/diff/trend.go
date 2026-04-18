package diff

import (
	"sort"
	"time"
)

// TrendPoint represents drift metrics at a point in time.
type TrendPoint struct {
	Timestamp   time.Time `json:"timestamp"`
	Total       int       `json:"total"`
	OnlySource  int       `json:"only_source"`
	OnlyDest    int       `json:"only_destination"`
	Modified    int       `json:"modified"`
}

// Trend holds an ordered series of TrendPoints.
type Trend struct {
	Points []TrendPoint `json:"points"`
}

// BuildTrend converts a slice of HistoryEntry values into a Trend.
// Entries are sorted by timestamp ascending before processing.
func BuildTrend(entries []HistoryEntry) Trend {
	sorted := make([]HistoryEntry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Timestamp.Before(sorted[j].Timestamp)
	})

	points := make([]TrendPoint, 0, len(sorted))
	for _, e := range sorted {
		var src, dst, mod int
		for _, r := range e.Results {
			switch r.Status {
			case "only_source":
				src++
			case "only_destination":
				dst++
			case "modified":
				mod++
			}
		}
		points = append(points, TrendPoint{
			Timestamp:  e.Timestamp,
			Total:      src + dst + mod,
			OnlySource: src,
			OnlyDest:   dst,
			Modified:   mod,
		})
	}
	return Trend{Points: points}
}

// Delta returns the change in total drift between the first and last point.
// Returns 0 if fewer than two points exist.
func (t Trend) Delta() int {
	if len(t.Points) < 2 {
		return 0
	}
	return t.Points[len(t.Points)-1].Total - t.Points[0].Total
}
