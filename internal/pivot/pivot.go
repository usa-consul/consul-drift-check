// Package pivot groups diff results into a key→datacenter matrix,
// making it easy to compare values for the same key across multiple DCs.
package pivot

import (
	"sort"

	"github.com/your-org/consul-drift-check/internal/diff"
)

// Row holds the values observed for a single KV key across datacenters.
type Row struct {
	Key    string
	Values map[string]string // dc label → value (empty string = absent)
}

// Table is an ordered slice of Rows.
type Table []Row

// Build constructs a pivot Table from per-DC diff results.
// dcResults maps a datacenter label to the slice of diff.Result produced
// by comparing that DC against the source.
func Build(dcResults map[string][]diff.Result) Table {
	type cell struct {
		dc    string
		value string
	}

	agg := make(map[string][]cell)

	for dc, results := range dcResults {
		for _, r := range results {
			agg[r.Key] = append(agg[r.Key], cell{dc: dc, value: r.DestinationValue})
		}
	}

	keys := make([]string, 0, len(agg))
	for k := range agg {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	table := make(Table, 0, len(keys))
	for _, k := range keys {
		row := Row{
			Key:    k,
			Values: make(map[string]string, len(agg[k])),
		}
		for _, c := range agg[k] {
			row.Values[c.dc] = c.value
		}
		table = append(table, row)
	}
	return table
}

// DCs returns a sorted list of all datacenter labels present in the table.
func (t Table) DCs() []string {
	seen := make(map[string]struct{})
	for _, row := range t {
		for dc := range row.Values {
			seen[dc] = struct{}{}
		}
	}
	dcs := make([]string, 0, len(seen))
	for dc := range seen {
		dcs = append(dcs, dc)
	}
	sort.Strings(dcs)
	return dcs
}
