// Package truncate limits KV pair values to a maximum byte length.
// This is useful when exporting or displaying values that may be very large.
package truncate

import (
	"github.com/hashicorp/consul/api"
)

// Options controls truncation behaviour.
type Options struct {
	// MaxBytes is the maximum number of bytes to retain per value.
	// Values longer than this are trimmed and a suffix is appended.
	// A value of 0 disables truncation.
	MaxBytes int

	// Suffix is appended to truncated values. Defaults to "..." when empty.
	Suffix string
}

// Apply truncates the Value field of each KV pair according to opts.
// The input slice is not modified; a new slice is returned.
func Apply(pairs []*api.KVPair, opts Options) []*api.KVPair {
	if opts.MaxBytes <= 0 {
		return pairs
	}
	suffix := opts.Suffix
	if suffix == "" {
		suffix = "..."
	}
	out := make([]*api.KVPair, len(pairs))
	for i, p := range pairs {
		copy := *p
		if len(copy.Value) > opts.MaxBytes {
			copy.Value = append([]byte(nil), copy.Value[:opts.MaxBytes]...)
			copy.Value = append(copy.Value, []byte(suffix)...)
		}
		out[i] = &copy
	}
	return out
}
