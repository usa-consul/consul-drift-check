// Package diff compares Consul KV pairs between a source and destination
// namespace and reports which keys are only in source, only in destination,
// or present in both but with differing values.
//
// It also provides helpers for persisting and querying historical drift runs
// via AppendHistory / LoadHistory / Since.
package diff
