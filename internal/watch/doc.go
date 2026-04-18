// Package watch provides key-level diffing between two sets of Consul KV
// pairs, identifying keys that were added, removed, or modified.
//
// It is intended to be used alongside the watcher package to detect
// incremental changes between polling intervals.
package watch
