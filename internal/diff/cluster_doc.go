// Package diff provides utilities for comparing Consul KV namespaces.
//
// The cluster sub-feature (cluster.go) aggregates per-key drift results
// into a ClusterSummary, grouping entries by their top-level path segment
// and annotating each group with the source/destination datacenter pair.
//
// Typical usage:
//
//	summary := diff.BuildCluster("dc1", "dc2", results)
//	for _, p := range summary.Pairs {
//		fmt.Printf("%s: +%d -%d ~%d\n", p.Pair, p.Added, p.Removed, p.Modified)
//	}
package diff
