package diff

import "sort"

// SliceOptions controls how result slicing is applied.
type SliceOptions struct {
	// Offset is the zero-based index of the first result to include.
	Offset int
	// Limit is the maximum number of results to return. Zero means no limit.
	Limit int
	// OrderByKey sorts results alphabetically before slicing.
	OrderByKey bool
}

// Slice returns a subset of results using offset/limit pagination.
// If opts.OrderByKey is true the slice is applied after sorting.
func Slice(results []Result, opts SliceOptions) []Result {
	if len(results) == 0 {
		return nil
	}

	out := make([]Result, len(results))
	copy(out, results)

	if opts.OrderByKey {
		sort.Slice(out, func(i, j int) bool {
			return out[i].Key < out[j].Key
		})
	}

	if opts.Offset >= len(out) {
		return nil
	}
	out = out[opts.Offset:]

	if opts.Limit > 0 && opts.Limit < len(out) {
		out = out[:opts.Limit]
	}

	return out
}
