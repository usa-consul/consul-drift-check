// Package diff compares Consul KV namespaces and reports configuration drift.
//
// Core functions:
//   - Compare: diff two KV maps and return a slice of Result
//   - Patch: convert results into KV set/delete operations
//   - Rollup: aggregate results by top-level key segment
//   - AppendHistory / LoadHistory / Since: persist and query drift history
//   - BuildTrend: derive a time-series Trend from history entries
package diff
