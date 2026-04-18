// Package prefix provides utilities for analysing and normalising KV key prefixes.
package prefix

import (
	"sort"
	"strings"

	"github.com/hashicorp/consul/api"
)

// Stat holds aggregated information about a single prefix.
type Stat struct {
	Prefix string
	Count  int
	Keys   []string
}

// Analyse groups the supplied KV pairs by their top-level path segment and
// returns one Stat per distinct prefix, sorted alphabetically.
func Analyse(pairs []*api.KVPair) []Stat {
	groups := make(map[string][]string)
	for _, p := range pairs {
		if p == nil {
			continue
		}
		seg := topSegment(p.Key)
		groups[seg] = append(groups[seg], p.Key)
	}

	stats := make([]Stat, 0, len(groups))
	for prefix, keys := range groups {
		sort.Strings(keys)
		stats = append(stats, Stat{
			Prefix: prefix,
			Count:  len(keys),
			Keys:   keys,
		})
	}
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Prefix < stats[j].Prefix
	})
	return stats
}

// Common returns prefixes that appear in both a and b.
func Common(a, b []*api.KVPair) []string {
	setA := prefixSet(a)
	var common []string
	for _, p := range b {
		if p == nil {
			continue
		}
		seg := topSegment(p.Key)
		if setA[seg] {
			common = append(common, seg)
			delete(setA, seg) // emit once
		}
	}
	sort.Strings(common)
	return common
}

func prefixSet(pairs []*api.KVPair) map[string]bool {
	m := make(map[string]bool)
	for _, p := range pairs {
		if p != nil {
			m[topSegment(p.Key)] = true
		}
	}
	return m
}

func topSegment(key string) string {
	key = strings.TrimPrefix(key, "/")
	if idx := strings.IndexByte(key, '/'); idx != -1 {
		return key[:idx]
	}
	return key
}
