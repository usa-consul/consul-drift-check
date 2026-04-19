// Package diff provides utilities for comparing Consul KV namespaces.
//
// # Budget
//
// EvaluateBudget checks a slice of diff Results against a BudgetOptions
// threshold. It returns a BudgetResult that reports how many keys drifted in
// each category and whether any configured limit was exceeded.
//
// Typical usage:
//
//	result := diff.EvaluateBudget(results, diff.BudgetOptions{
//		MaxDriftCount: 50,
//		MaxModified:   20,
//		MaxMissing:    10,
//	})
//	if result.Exceeded {
//		log.Println("drift budget exceeded:", result.Violations)
//	}
package diff
