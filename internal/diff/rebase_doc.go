// Package diff provides utilities for comparing and transforming Consul KV
// diff results.
//
// # Rebase
//
// Rebase rewrites the key paths of a []Result slice so that results belonging
// to one namespace can be mapped onto a different prefix root before further
// processing or reporting.
//
// Example:
//
//	results := Rebase(raw, RebaseOptions{
//		StripPrefix:   "prod/service-a",
//		NewPrefix:     "staging/service-a",
//		SkipUnmatched: true,
//	})
package diff
