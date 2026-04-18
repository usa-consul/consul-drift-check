// Package pivot transforms per-datacenter diff results into a
// key-centric matrix (Table) where each row represents one KV key
// and each column represents the value observed in a specific datacenter.
//
// This view is useful when drift is being checked across more than two
// datacenters simultaneously and the operator wants to compare all values
// side-by-side rather than reading individual pairwise reports.
package pivot
