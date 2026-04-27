// Package diff provides utilities for comparing Consul KV namespaces.
//
// # Supersede
//
// Supersede removes drift results that are shadowed by higher-priority prefix
// entries sharing the same top-level key segment.
//
// Use case: when a canonical "config/app" prefix is the authoritative source
// of truth, any drift detected under a legacy "config/" path that does NOT
// fall under "config/app" can be suppressed so that operators focus only on
// the entries that matter.
//
// Example:
//
//	results := Supersede(all, SupersedeOptions{
//	    Prefixes: []string{"config/app"},
//	})
package diff
