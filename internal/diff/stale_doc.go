// Package diff provides utilities for comparing Consul KV namespaces and
// analysing the resulting drift results.
//
// # Stale Detection
//
// DetectStale annotates each Result with an Age and an IsStale flag.
// A result is considered stale when the duration since the key was first
// observed (according to the caller-supplied seenAt map) meets or exceeds
// StaleOptions.MaxAge.
//
// Typical usage:
//
//	sr := diff.DetectStale(results, seenAt, time.Now(), diff.StaleOptions{
//		MaxAge:     24 * time.Hour,
//		OnlyStatus: "modified",
//	})
package diff
