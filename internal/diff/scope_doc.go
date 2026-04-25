// Package diff provides utilities for comparing Consul KV namespaces and
// analysing the resulting drift results.
//
// Scope filters a result set to a defined operational boundary by combining
// prefix matching, status filtering, and an optional hard cap on the number
// of returned entries. It is useful when an operator wants to inspect drift
// within a specific subtree (e.g. "service/payments") or for a particular
// change category (e.g. only "modified" keys) without processing the full
// result set.
package diff
