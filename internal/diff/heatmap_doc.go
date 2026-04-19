// Package diff provides utilities for comparing Consul KV namespaces and
// analysing drift over time.
//
// # HeatMap
//
// BuildHeatMap aggregates drift history entries into a ranked heatmap that
// surfaces which top-level key prefixes experience the most configuration
// drift. Each HeatCell records the total number of drifted keys and a
// weighted severity score (modified keys score 2.0; missing keys score 1.0).
//
// Typical usage:
//
//	entries, _ := diff.LoadHistory("drift-history.ndjson")
//	hm := diff.BuildHeatMap(entries)
//	for _, cell := range hm {
//		fmt.Printf("%s\tscore=%.1f\tcount=%d\n", cell.Prefix, cell.Score, cell.Count)
//	}
package diff
