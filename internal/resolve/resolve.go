// Package resolve provides key resolution utilities that map raw Consul KV
// keys to human-readable aliases using a configurable rule set.
package resolve

import (
	"strings"

	"github.com/hashicorp/consul/api"
)

// Rule maps a key prefix to a human-readable alias.
type Rule struct {
	Prefix string `yaml:"prefix"`
	Alias  string `yaml:"alias"`
}

// Options controls how resolution is applied.
type Options struct {
	Rules         []Rule
	FallbackToKey bool
}

// Result holds a KV pair with its resolved alias.
type Result struct {
	Pair  *api.KVPair
	Alias string
}

// Apply resolves each KV pair to an alias based on the configured rules.
// If no rule matches and FallbackToKey is true, the raw key is used as the alias.
func Apply(pairs []*api.KVPair, opts Options) []Result {
	results := make([]Result, 0, len(pairs))
	for _, p := range pairs {
		if p == nil {
			continue
		}
		alias := resolve(p.Key, opts)
		results = append(results, Result{Pair: p, Alias: alias})
	}
	return results
}

func resolve(key string, opts Options) string {
	for _, r := range opts.Rules {
		prefix := strings.TrimSuffix(r.Prefix, "/")
		if strings.HasPrefix(key, prefix) {
			suffix := strings.TrimPrefix(key, prefix)
			suffix = strings.TrimPrefix(suffix, "/")
			if suffix == "" {
				return r.Alias
			}
			return r.Alias + "/" + suffix
		}
	}
	if opts.FallbackToKey {
		return key
	}
	return ""
}
