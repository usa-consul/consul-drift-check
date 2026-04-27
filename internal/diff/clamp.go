package diff

// ClampOptions controls how value lengths are constrained in diff results.
type ClampOptions struct {
	// MaxValueLen is the maximum number of bytes kept for SourceValue and
	// DestValue. Zero means no clamping.
	MaxValueLen int
	// Suffix is appended when a value is truncated. Defaults to "…".
	Suffix string
}

// Clamp truncates the SourceValue and DestValue fields of every Result so that
// neither exceeds MaxValueLen bytes. Results whose values are already within
// the limit are returned unchanged. The original slice is never mutated; a new
// slice is always returned.
//
// If opts.MaxValueLen is zero or negative Clamp returns a sorted copy of the
// input without any truncation.
func Clamp(results []Result, opts ClampOptions) []Result {
	out := make([]Result, len(results))
	copy(out, results)
	sortByKey(out)

	if opts.MaxValueLen <= 0 {
		return out
	}

	suffix := opts.Suffix
	if suffix == "" {
		suffix = "\u2026" // …
	}

	for i, r := range out {
		out[i].SourceValue = clampBytes(r.SourceValue, opts.MaxValueLen, suffix)
		out[i].DestValue = clampBytes(r.DestValue, opts.MaxValueLen, suffix)
	}
	return out
}

func clampBytes(b []byte, max int, suffix string) []byte {
	if len(b) <= max {
		return b
	}
	suf := []byte(suffix)
	keep := max - len(suf)
	if keep < 0 {
		keep = 0
	}
	out := make([]byte, keep+len(suf))
	copy(out, b[:keep])
	copy(out[keep:], suf)
	return out
}
