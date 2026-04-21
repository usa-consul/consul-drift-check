package diff

import (
	"sort"
	"strings"
)

// RebaseOptions controls how results are rebased onto a new prefix root.
type RebaseOptions struct {
	// StripPrefix removes this prefix from each key before applying NewPrefix.
	StripPrefix string
	// NewPrefix is prepended to every key after stripping.
	NewPrefix string
	// SkipUnmatched drops results whose key does not start with StripPrefix.
	SkipUnmatched bool
}

// Rebase rewrites the Key field of each Result according to opts, returning a
// new slice with keys normalised to the requested prefix root. The original
// slice is never modified.
func Rebase(results []Result, opts RebaseOptions) []Result {
	if len(results) == 0 {
		return nil
	}

	strip := normaliseRebasePrefix(opts.StripPrefix)
	newPfx := normaliseRebasePrefix(opts.NewPrefix)

	out := make([]Result, 0, len(results))
	for _, r := range results {
		key := r.Key

		if strip != "" {
			if !strings.HasPrefix(key, strip) {
				if opts.SkipUnmatched {
					continue
				}
			} else {
				key = strings.TrimPrefix(key, strip)
			}
		}

		if newPfx != "" {
			key = newPfx + key
		}

		clone := r
		clone.Key = key
		out = append(out, clone)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Key < out[j].Key
	})
	return out
}

func normaliseRebasePrefix(p string) string {
	p = strings.TrimSpace(p)
	if p == "" {
		return ""
	}
	if !strings.HasSuffix(p, "/") {
		p += "/"
	}
	return p
}
