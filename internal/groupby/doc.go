// Package groupby partitions a slice of diff.Result values into named
// groups based on a configurable path segment index.
//
// This is useful when reporting drift summaries per service, environment,
// or any other logical namespace encoded in the KV key hierarchy.
package groupby
