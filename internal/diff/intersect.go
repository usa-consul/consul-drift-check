package diff

import "sort"

// IntersectResult holds a key that exists in both source and destination
// along with its values from each side.
type IntersectResult struct {
	Key   string
	Src   []byte
	Dst   []byte
	Equal bool
}

// Intersect returns only the keys present in both src and dst maps,
// annotating each with whether the values are equal.
func Intersect(src, dst map[string][]byte) []IntersectResult {
	if len(src) == 0 || len(dst) == 0 {
		return nil
	}

	var results []IntersectResult

	for k, sv := range src {
		dv, ok := dst[k]
		if !ok {
			continue
		}
		results = append(results, IntersectResult{
			Key:   k,
			Src:   sv,
			Dst:   dv,
			Equal: bytesEqual(sv, dv),
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})

	return results
}

// IntersectKeys returns the sorted list of keys common to both maps.
func IntersectKeys(src, dst map[string][]byte) []string {
	results := Intersect(src, dst)
	keys := make([]string, len(results))
	for i, r := range results {
		keys[i] = r.Key
	}
	return keys
}
