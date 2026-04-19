package diff

import (
	"sort"
	"strings"
)

// HeatCell represents drift activity for a single key prefix.
type HeatCell struct {
	Prefix string
	Count  int
	Score  float64
}

// HeatMap is an ordered slice of HeatCells, highest score first.
type HeatMap []HeatCell

// BuildHeatMap aggregates drift results across multiple history entries
// into a ranked heatmap keyed by top-level prefix segment.
func BuildHeatMap(entries []HistoryEntry) HeatMap {
	type acc struct {
		count int
		score float64
	}
	buckets := make(map[string]*acc)

	for _, e := range entries {
		for _, r := range e.Results {
			seg := topHeatSegment(r.Key)
			if _, ok := buckets[seg]; !ok {
				buckets[seg] = &acc{}
			}
			buckets[seg].count++
			buckets[seg].score += severityScore(r.Status)
		}
	}

	hm := make(HeatMap, 0, len(buckets))
	for prefix, a := range buckets {
		hm = append(hm, HeatCell{
			Prefix: prefix,
			Count:  a.count,
			Score:  a.score,
		})
	}

	sort.Slice(hm, func(i, j int) bool {
		if hm[i].Score != hm[j].Score {
			return hm[i].Score > hm[j].Score
		}
		return hm[i].Prefix < hm[j].Prefix
	})
	return hm
}

func topHeatSegment(key string) string {
	key = strings.TrimPrefix(key, "/")
	if idx := strings.Index(key, "/"); idx >= 0 {
		return key[:idx]
	}
	return key
}

func severityScore(status string) float64 {
	switch status {
	case "modified":
		return 2.0
	case "only_in_source", "only_in_destination":
		return 1.0
	}
	return 0.0
}
