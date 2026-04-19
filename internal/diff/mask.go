package diff

import (
	"strings"
)

// MaskOptions controls which result fields are masked in output.
type MaskOptions struct {
	// SensitivePrefixes hides values for keys whose path starts with any entry.
	SensitivePrefixes []string
	// Mask is the replacement string; defaults to "***".
	Mask string
}

// MaskedResult is a Result with potentially redacted values.
type MaskedResult struct {
	Result
	SourceValue string
	DestValue   string
}

// Mask applies value redaction to a slice of Results based on MaskOptions.
// Keys matching a sensitive prefix have their source and destination values
// replaced with the configured mask string.
func Mask(results []Result, opts MaskOptions) []MaskedResult {
	mask := opts.Mask
	if mask == "" {
		mask = "***"
	}

	out := make([]MaskedResult, 0, len(results))
	for _, r := range results {
		mr := MaskedResult{
			Result:      r,
			SourceValue: string(r.SourceValue),
			DestValue:   string(r.DestValue),
		}
		if isMaskedKey(r.Key, opts.SensitivePrefixes) {
			if mr.SourceValue != "" {
				mr.SourceValue = mask
			}
			if mr.DestValue != "" {
				mr.DestValue = mask
			}
		}
		out = append(out, mr)
	}
	return out
}

func isMaskedKey(key string, prefixes []string) bool {
	norm := strings.ToLower(key)
	for _, p := range prefixes {
		if strings.HasPrefix(norm, strings.ToLower(p)) {
			return true
		}
	}
	return false
}
