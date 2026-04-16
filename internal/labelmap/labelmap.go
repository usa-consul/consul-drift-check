// Package labelmap maps Consul KV keys to structured labels for grouping and reporting.
package labelmap

import (
	"strings"

	"github.com/hashicorp/consul/api"
)

// Rule maps a key prefix to a set of labels.
type Rule struct {
	Prefix string            `yaml:"prefix"`
	Labels map[string]string `yaml:"labels"`
}

// Result holds a KV pair with its resolved labels.
type Result struct {
	Key    string
	Value  []byte
	Labels map[string]string
}

// Apply resolves labels for each KV pair using the provided rules.
// The first matching rule wins. Pairs with no match receive an empty label map.
func Apply(pairs []*api.KVPair, rules []Rule) []Result {
	results := make([]Result, 0, len(pairs))
	for _, p := range pairs {
		results = append(results, Result{
			Key:    p.Key,
			Value:  p.Value,
			Labels: resolveLabels(p.Key, rules),
		})
	}
	return results
}

func resolveLabels(key string, rules []Rule) map[string]string {
	for _, r := range rules {
		if strings.HasPrefix(key, r.Prefix) {
			out := make(map[string]string, len(r.Labels))
			for k, v := range r.Labels {
				out[k] = v
			}
			return out
		}
	}
	return map[string]string{}
}
