// Package tag attaches arbitrary string tags to diff results based on key prefix rules.
package tag

import (
	"strings"

	"github.com/your-org/consul-drift-check/internal/diff"
)

// Rule maps a key prefix to a set of tags.
type Rule struct {
	Prefix string
	Tags   []string
}

// Options controls tagging behaviour.
type Options struct {
	Rules []Rule
}

// TaggedResult wraps a diff result with resolved tags.
type TaggedResult struct {
	diff.Result
	Tags []string
}

// Apply resolves tags for each result according to the provided rules.
// The first matching rule wins. Results with no matching rule receive an empty tag list.
func Apply(results []diff.Result, opts Options) []TaggedResult {
	out := make([]TaggedResult, 0, len(results))
	for _, r := range results {
		out = append(out, TaggedResult{
			Result: r,
			Tags:   resolveTags(r.Key, opts.Rules),
		})
	}
	return out
}

func resolveTags(key string, rules []Rule) []string {
	for _, rule := range rules {
		if strings.HasPrefix(key, rule.Prefix) {
			copied := make([]string, len(rule.Tags))
			copy(copied, rule.Tags)
			return copied
		}
	}
	return []string{}
}
