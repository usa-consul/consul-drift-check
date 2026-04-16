// Package normalize provides utilities for standardising Consul KV keys
// before comparison, ensuring consistent formatting across datacenters.
package normalize

import (
	"strings"

	"github.com/hashicorp/consul/api"
)

// Options controls how keys are normalised.
type Options struct {
	// LowerCase converts all keys to lower case.
	LowerCase bool
	// StripPrefix removes a leading prefix from every key.
	StripPrefix string
	// CollapseSlashes replaces consecutive slashes with a single slash.
	CollapseSlashes bool
}

// Apply returns a new slice of KV pairs with keys normalised according to opts.
// Values are never modified.
func Apply(pairs []*api.KVPair, opts Options) []*api.KVPair {
	out := make([]*api.KVPair, 0, len(pairs))
	for _, p := range pairs {
		key := p.Key
		key = applyStripPrefix(key, opts.StripPrefix)
		if opts.CollapseSlashes {
			key = collapseSlashes(key)
		}
		if opts.LowerCase {
			key = strings.ToLower(key)
		}
		out = append(out, &api.KVPair{
			Key:   key,
			Value: p.Value,
			Flags: p.Flags,
		})
	}
	return out
}

func applyStripPrefix(key, prefix string) string {
	if prefix == "" {
		return key
	}
	p := strings.TrimSuffix(prefix, "/") + "/"
	return strings.TrimPrefix(key, p)
}

func collapseSlashes(key string) string {
	var b strings.Builder
	prev := rune(0)
	for _, ch := range key {
		if ch == '/' && prev == '/' {
			continue
		}
		b.WriteRune(ch)
		prev = ch
	}
	return b.String()
}
