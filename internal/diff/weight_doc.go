// Package diff provides utilities for comparing Consul KV namespaces.
//
// # Weight
//
// The weight sub-feature assigns a numeric importance score to each drift
// result so that operators can prioritise which discrepancies to resolve first.
//
// Base scores are configurable per status (modified, only-in-source,
// only-in-destination). An optional map of key-prefix multipliers allows
// critical namespaces (e.g. "service/") to receive elevated scores.
//
// Results returned by [Apply] are sorted in descending order of weight so the
// highest-priority drifts appear first.
package diff
