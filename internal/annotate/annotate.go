// Package annotate attaches free-form metadata annotations to diff results.
package annotate

import (
	"strings"

	"github.com/example/consul-drift-check/internal/diff"
)

// Rule maps a key prefix to a set of annotations.
type Rule struct {
	Prefix      string            `yaml:"prefix"`
	Annotations map[string]string `yaml:"annotations"`
}

// Result wraps a diff result with additional annotations.
type Result struct {
	diff.Result
	Annotations map[string]string
}

// Apply attaches annotations to each diff result according to the provided
// rules. The first matching rule wins. Results with no matching rule receive
// an empty annotation map.
func Apply(results []diff.Result, rules []Rule) []Result {
	out := make([]Result, 0, len(results))
	for _, r := range results {
		out = append(out, Result{
			Result:      r,
			Annotations: resolve(r.Key, rules),
		})
	}
	return out
}

func resolve(key string, rules []Rule) map[string]string {
	for _, rule := range rules {
		if strings.HasPrefix(key, rule.Prefix) {
			copy := make(map[string]string, len(rule.Annotations))
			for k, v := range rule.Annotations {
				copy[k] = v
			}
			return copy
		}
	}
	return map[string]string{}
}
