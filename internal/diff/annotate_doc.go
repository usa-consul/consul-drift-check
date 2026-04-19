// Package diff — annotate.go
//
// Annotate enriches a slice of diff Results with human-readable labels
// derived from key-prefix rules. This is useful when rendering reports
// or exporting data to external systems that require a classification
// field alongside each drift entry.
//
// Usage:
//
//	annotated := diff.Annotate(results, diff.AnnotateOptions{
//		PrefixLabels: map[string]string{
//			"service/": "services",
//			"infra/":   "infrastructure",
//		},
//		DefaultLabel: "general",
//	})
package diff
