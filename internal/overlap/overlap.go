// Package overlap identifies keys that appear in more than one namespace prefix.
package overlap

import (
	"sort"
	"strings"

	"github.com/hashicorp/consul/api"
)

// Result holds a key that overlaps across two or more prefixes.
type Result struct {
	Key      string   `json:"key"`
	Prefixes []string `json:"prefixes"`
}

// Find returns keys from pairs that match more than one of the given prefixes.
// Each Result lists the prefixes the key falls under.
func Find(pairs []*api.KVPair, prefixes []string) []Result {
	if len(pairs) == 0 || len(prefixes) == 0 {
		return nil
	}

	type entry struct {
		matches map[string]struct{}
	}

	index := make(map[string]*entry, len(pairs))
	for _, p := range pairs {
		if p == nil {
			continue
		}
		e := &entry{matches: make(map[string]struct{})}
		for _, pfx := range prefixes {
			norm := strings.TrimSuffix(pfx, "/") + "/"
			if strings.HasPrefix(p.Key, norm) || p.Key == strings.TrimSuffix(pfx, "/") {
				e.matches[pfx] = struct{}{}
			}
		}
		if len(e.matches) > 1 {
			index[p.Key] = e
		}
	}

	results := make([]Result, 0, len(index))
	for key, e := range index {
		pfxList := make([]string, 0, len(e.matches))
		for pfx := range e.matches {
			pfxList = append(pfxList, pfx)
		}
		sort.Strings(pfxList)
		results = append(results, Result{Key: key, Prefixes: pfxList})
	}
	sort.Slice(results, func(i, j int) bool { return results[i].Key < results[j].Key })
	return results
}

// Keys returns only the key strings from a slice of Results.
func Keys(results []Result) []string {
	out := make([]string, len(results))
	for i, r := range results {
		out[i] = r.Key
	}
	return out
}
