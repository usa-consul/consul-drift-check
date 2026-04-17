// Package rank orders diff results by severity and drift magnitude.
package rank

import (
	"sort"

	"github.com/example/consul-drift-check/internal/diff"
)

// Order defines how results are sorted.
type Order string

const (
	OrderSeverity  Order = "severity"
	OrderKey       Order = "key"
	OrderStatus    Order = "status"
)

var statusWeight = map[diff.Status]int{
	diff.StatusModified:    3,
	diff.StatusOnlyInSrc:   2,
	diff.StatusOnlyInDst:   2,
	diff.StatusMatch:       0,
}

// Apply returns a new slice of results sorted by the given order.
// The original slice is not modified.
func Apply(results []diff.Result, order Order) []diff.Result {
	out := make([]diff.Result, len(results))
	copy(out, results)

	switch order {
	case OrderKey:
		sort.Slice(out, func(i, j int) bool {
			return out[i].Key < out[j].Key
		})
	case OrderStatus:
		sort.Slice(out, func(i, j int) bool {
			return out[i].Status < out[j].Status
		})
	case OrderSeverity:
		fallthrough
	default:
		sort.Slice(out, func(i, j int) bool {
			wi := statusWeight[out[i].Status]
			wj := statusWeight[out[j].Status]
			if wi != wj {
				return wi > wj
			}
			return out[i].Key < out[j].Key
		})
	}

	return out
}
