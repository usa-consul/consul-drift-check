// Package truncate provides value-length limiting for Consul KV pairs.
//
// Large values can cause issues when rendering reports or exporting data.
// Apply returns a new slice of KV pairs with values trimmed to a configurable
// maximum byte length. A customisable suffix (default "...") is appended to
// any value that was shortened so consumers can detect truncation.
package truncate
