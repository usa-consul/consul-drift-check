package diff

import (
	"strings"
	"time"
)

// StaleOptions configures stale-entry detection.
type StaleOptions struct {
	// MaxAge is the duration after which an entry is considered stale.
	// Entries whose ModifyIndex has not changed within MaxAge are flagged.
	MaxAge time.Duration

	// OnlyStatus restricts staleness checks to results with the given status.
	// Leave empty to check all statuses.
	OnlyStatus string

	// ExcludePrefixes lists key prefixes that are exempt from staleness checks.
	ExcludePrefixes []string
}

// StaleResult wraps a Result with staleness metadata.
type StaleResult struct {
	Result
	Age     time.Duration
	IsStale bool
}

// DetectStale annotates results that have not changed within opts.MaxAge.
// The seenAt map provides the last-observed timestamp for each key; keys
// absent from the map are treated as first-seen at referenceTime.
func DetectStale(results []Result, seenAt map[string]time.Time, referenceTime time.Time, opts StaleOptions) []StaleResult {
	if len(results) == 0 {
		return nil
	}

	out := make([]StaleResult, 0, len(results))

	for _, r := range results {
		sr := StaleResult{Result: r}

		if opts.OnlyStatus != "" && !strings.EqualFold(r.Status, opts.OnlyStatus) {
			out = append(out, sr)
			continue
		}

		if isStaleExcluded(r.Key, opts.ExcludePrefixes) {
			out = append(out, sr)
			continue
		}

		if opts.MaxAge <= 0 {
			out = append(out, sr)
			continue
		}

		first, ok := seenAt[r.Key]
		if !ok {
			first = referenceTime
		}

		age := referenceTime.Sub(first)
		sr.Age = age
		sr.IsStale = age >= opts.MaxAge
		out = append(out, sr)
	}

	return out
}

func isStaleExcluded(key string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}
