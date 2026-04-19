// Package diff provides utilities for comparing Consul KV namespaces.
//
// # Slice
//
// Slice applies offset/limit pagination to a []Result slice. It is useful
// when rendering paginated output or when only a subset of drift results
// needs to be processed downstream.
//
// Results can optionally be sorted by key before slicing to ensure
// deterministic pages regardless of the order in which they were produced.
package diff
