// Package sanitize normalises Consul KV keys before comparison.
package sanitize

import (
	"strings"

	"github.com/hashicorp/consul/api"
)

// Options controls how keys are normalised.
type Options struct {
	// LowerCase converts all keys to lower-case before comparison.
	LowerCase bool
	// TrimPrefix removes a leading path segment from every key.
	TrimPrefix string
}

// Apply returns a new slice of KV pairs with keys normalised according to opts.
// The original slice is not modified.
func Apply(pairs api.KVPairs, opts Options) api.KVPairs {
	out := make(api.KVPairs, 0, len(pairs))
	for _, p := range pairs {
		if p == nil {
			continue
		}
		copy := *p
		copy.Key = normalise(p.Key, opts)
		out = append(out, &copy)
	}
	return out
}

func normalise(key string, opts Options) string {
	if opts.TrimPrefix != "" {
		prefix := opts.TrimPrefix
		if !strings.HasSuffix(prefix, "/") {
			prefix += "/"
		}
		key = strings.TrimPrefix(key, prefix)
	}
	if opts.LowerCase {
		key = strings.ToLower(key)
	}
	return key
}
