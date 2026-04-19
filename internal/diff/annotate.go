// Package diff provides utilities for comparing Consul KV namespaces.
package diff

import (
	"strings"
)

// AnnotateOptions controls how results are annotated.
type AnnotateOptions struct {
	// PrefixLabels maps key prefixes to human-readable labels.
	PrefixLabels map[string]string
	// DefaultLabel is used when no prefix matches.
	DefaultLabel string

// AnnotatedResult wraps a Result with an additional label.
type AnnotatedResult struct {
	Result
	Label Annotate attaches a label to each Result based on its key prefix.
// The first matching prefix wins.Label is used.
func Annotate(results []Result, opts AnnotateOptions) []AnnotatedResult {
	if len(results) == 0 {
		return nil
	}
	out := make([]AnnotatedResult, 0, len(results))
	for _, r := range results {
		out = append(out, AnnotatedResult{
			Result: r,
			Label:  resolveAnnotateLabel(r.Key, opts),
		})
	}
	return out
}

func resolveAnnotateLabel(key string, opts AnnotateOptions) string {
	for prefix, label := range opts.PrefixLabels {
		p := strings.TrimRight(prefix, "/") + "/"
		if strings.HasPrefix(key, p) || key == strings.TrimRight(prefix, "/") {
			return label
		}
	}
	return opts.DefaultLabel
}
