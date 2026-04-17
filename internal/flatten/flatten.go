// Package flatten converts nested map structures from Consul KV into
// a flat key/value representation suitable for drift comparison.
package flatten

import (
	"strings"

	"github.com/hashicorp/consul/api"
)

// Options controls how flattening is applied.
type Options struct {
	// Separator is placed between path segments when joining nested keys.
	// Defaults to "/" when empty.
	Separator string
	// StripPrefix removes a leading prefix from every key before storing.
	StripPrefix string
}

// Apply takes a slice of Consul KV pairs and returns a flat map of
// key -> raw value bytes. Pairs with nil values are kept with an empty
// byte slice so that missing-key drift is still detectable.
func Apply(pairs api.KVPairs, opts Options) map[string][]byte {
	sep := opts.Separator
	if sep == "" {
		sep = "/"
	}

	prefix := normalisePrefix(opts.StripPrefix, sep)

	out := make(map[string][]byte, len(pairs))
	for _, p := range pairs {
		if p == nil {
			continue
		}
		k := p.Key
		if prefix != "" {
			k = strings.TrimPrefix(k, prefix)
		}
		k = strings.Trim(k, sep)
		if k == "" {
			continue
		}
		v := p.Value
		if v == nil {
			v = []byte{}
		}
		out[k] = v
	}
	return out
}

func normalisePrefix(p, sep string) string {
	if p == "" {
		return ""
	}
	p = strings.TrimRight(p, sep)
	return p + sep
}
