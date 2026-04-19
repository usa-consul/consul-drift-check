package diff

import (
	"sort"
	"strings"
)

// ClusterResult holds drift counts for a named datacenter pair.
type ClusterResult struct {
	Pair     string
	Source   string
	Dest     string
	Added    int
	Removed  int
	Modified int
	Total    int
}

// ClusterSummary aggregates drift results across multiple DC pairs.
type ClusterSummary struct {
	Pairs []ClusterResult
}

// BuildCluster groups a slice of Result by datacenter pair label and
// returns a ClusterSummary sorted by pair name.
func BuildCluster(source, dest string, results []Result) ClusterSummary {
	counts := map[string]*ClusterResult{}

	for _, r := range results {
		seg := topClusterSegment(r.Key)
		if _, ok := counts[seg]; !ok {
			counts[seg] = &ClusterResult{
				Pair:   pairLabel(source, dest),
				Source: source,
				Dest:   dest,
			}
		}
		cr := counts[seg]
		switch r.Status {
		case StatusOnlyInSource:
			cr.Added++
		case StatusOnlyInDestination:
			cr.Removed++
		case StatusModified:
			cr.Modified++
		}
		cr.Total++
	}

	out := make([]ClusterResult, 0, len(counts))
	for _, v := range counts {
		out = append(out, *v)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Pair < out[j].Pair
	})
	return ClusterSummary{Pairs: out}
}

func pairLabel(src, dst string) string {
	return src + "->" + dst
}

func topClusterSegment(key string) string {
	key = strings.TrimPrefix(key, "/")
	if idx := strings.Index(key, "/"); idx >= 0 {
		return key[:idx]
	}
	return key
}
