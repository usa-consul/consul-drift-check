// Package diff provides utilities for comparing, analysing and forecasting
// Consul KV configuration drift between namespaces and datacenters.
//
// The Forecast function applies ordinary least-squares linear regression to
// a slice of TrendPoints and projects drift counts into the future. It is
// intended as a lightweight signal — not a replacement for proper time-series
// analysis — and works best when backed by at least a week of history.
package diff
