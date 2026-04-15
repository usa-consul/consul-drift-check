// Package baseline provides functionality for saving and loading a reference
// snapshot of Consul KV entries. A baseline captures the state of a prefix at
// a point in time and can later be compared against a live or current set of
// entries to detect configuration drift.
//
// Usage:
//
//	// Save a baseline from current source entries
//	err := baseline.Save("/var/lib/consul-drift/baseline.json", "app/", entries)
//
//	// Load and compare against current state
//	b, err := baseline.Load("/var/lib/consul-drift/baseline.json")
//	results := baseline.Compare(b, currentEntries)
package baseline
