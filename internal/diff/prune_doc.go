// Package diff provides utilities for comparing Consul KV namespaces.
//
// # Prune
//
// Prune removes unwanted entries from a diff result set before reporting.
// It supports three independent filter modes that can be combined:
//
//   - DropEmptyValues: discards keys whose value is empty in both source and
//     destination, which typically indicates placeholder keys with no real data.
//
//   - OnlyStatus: retains only results whose derived status (source_only,
//     dest_only, or modified) matches one of the supplied values. Useful for
//     focusing a report on additions or deletions only.
//
//   - StalePrefixes: drops any key that matches one of the supplied prefixes.
//     Intended for known-noisy namespaces that should be excluded from drift
//     alerting.
package diff
