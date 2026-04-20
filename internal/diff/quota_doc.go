// Package diff provides utilities for comparing Consul KV namespaces.
//
// # Quota
//
// EvaluateQuota checks whether the number of drifted keys under each
// top-level prefix exceeds a configured maximum. This allows operators
// to define acceptable drift budgets on a per-namespace basis.
//
// Example:
//
//	opts := diff.QuotaOptions{
//	    Rules:      map[string]int{"app": 5, "infra": 2},
//	    DefaultMax: 10,
//	}
//	violations := diff.EvaluateQuota(results, opts)
package diff
