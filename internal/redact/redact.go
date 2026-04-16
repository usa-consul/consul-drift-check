// Package redact masks sensitive values in KV pairs before output or logging.
package redact

import "strings"

// Options controls redaction behaviour.
type Options struct {
	// Keys whose values will be replaced with the Mask string.
	SensitiveKeys []string
	// Mask is the replacement string. Defaults to "***".
	Mask string
}

// KVPair is a minimal key/value pair used across the project.
type KVPair struct {
	Key   string
	Value []byte
}

const defaultMask = "***"

// Apply returns a new slice of KVPairs with sensitive values masked.
func Apply(pairs []KVPair, opts Options) []KVPair {
	mask := opts.Mask
	if mask == "" {
		mask = defaultMask
	}

	sensitive := make(map[string]struct{}, len(opts.SensitiveKeys))
	for _, k := range opts.SensitiveKeys {
		sensitive[normalise(k)] = struct{}{}
	}

	out := make([]KVPair, len(pairs))
	for i, p := range pairs {
		out[i] = p
		if isSensitive(p.Key, sensitive) {
			out[i].Value = []byte(mask)
		}
	}
	return out
}

func isSensitive(key string, sensitive map[string]struct{}) bool {
	nk := normalise(key)
	if _, ok := sensitive[nk]; ok {
		return true
	}
	// match by suffix segment, e.g. "password" matches "db/password"
	parts := strings.Split(nk, "/")
	_, ok := sensitive[parts[len(parts)-1]]
	return ok
}

func normalise(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
