// Package diff provides utilities for comparing Consul KV namespaces.
//
// # Normalize
//
// NormalizeResults and NormalizePairs transform key names before a diff is
// performed so that superficial differences in path layout (e.g. environment
// prefixes or casing conventions) do not generate false-positive drift
// results.
//
// Typical usage:
//
//	opts := diff.NormalizeOptions{StripPrefix: "prod", LowerCase: true}
//	src = diff.NormalizePairs(src, opts)
//	dst = diff.NormalizePairs(dst, opts)
//	results := diff.Compare(src, dst)
package diff
