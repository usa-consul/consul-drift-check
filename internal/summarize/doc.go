// Package summarize aggregates diff.Result slices into a structured Summary
// that exposes counts per change type (added, removed, modified) and a
// human-readable line-by-line report.
//
// Typical usage:
//
//	results := diff.Compare(src, dst)
//	summary := summarize.Build(results)
//	fmt.Println(summary.String())
//	fmt.Println(summary.Report())
package summarize
