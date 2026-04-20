package diff

import (
	"strings"

	"github.com/example/consul-drift-check/internal/consul"
)

// NormalizeOptions controls how result keys are normalised before comparison.
type NormalizeOptions struct {
	// StripPrefix removes a leading path segment from each result key.
	StripPrefix string
	// LowerCase converts all keys to lower case.
	LowerCase bool
}

// NormalizeResult returns a copy of r with the key transformed according to opts.
func NormalizeResult(r Result, opts NormalizeOptions) Result {
	key := r.Key
	if opts.StripPrefix != "" {
		prefix := strings.TrimSuffix(opts.StripPrefix, "/") + "/"
		key = strings.TrimPrefix(key, prefix)
	}
	if opts.LowerCase {
		key = strings.ToLower(key)
	}
	return Result{
		Key:         key,
		Status:      r.Status,
		SourceValue: r.SourceValue,
		DestValue:   r.DestValue,
	}
}

// NormalizeResults applies NormalizeResult to every element of results.
func NormalizeResults(results []Result, opts NormalizeOptions) []Result {
	if len(results) == 0 {
		return nil
	}
	out := make([]Result, len(results))
	for i, r := range results {
		out[i] = NormalizeResult(r, opts)
	}
	return out
}

// NormalizePairs applies the same key transformations to a slice of KV pairs.
func NormalizePairs(pairs []*consul.KVPair, opts NormalizeOptions) []*consul.KVPair {
	if len(pairs) == 0 {
		return nil
	}
	out := make([]*consul.KVPair, 0, len(pairs))
	for _, p := range pairs {
		if p == nil {
			continue
		}
		key := p.Key
		if opts.StripPrefix != "" {
			prefix := strings.TrimSuffix(opts.StripPrefix, "/") + "/"
			key = strings.TrimPrefix(key, prefix)
		}
		if opts.LowerCase {
			key = strings.ToLower(key)
		}
		out = append(out, &consul.KVPair{Key: key, Value: p.Value})
	}
	return out
}
