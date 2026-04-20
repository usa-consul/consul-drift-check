// Package diff provides utilities for comparing Consul KV namespaces.
//
// # Ceiling
//
// Ceiling limits the number of drift results returned per top-level KV prefix
// segment. This is useful when a single noisy prefix would otherwise dominate
// a report and obscure drift in other areas of the key space.
//
// Example:
//
//	results := diff.Ceiling(raw, diff.CeilingOptions{
//		MaxPerPrefix: 5,
//		OrderByKey:   true,
//	})
package diff
