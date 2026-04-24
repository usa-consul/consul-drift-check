package diff

import (
	"sort"
	"strings"
)

// PinOptions controls which keys are pinned and how pinning is enforced.
type PinOptions struct {
	// PinnedKeys is a list of exact keys whose values must not drift.
	PinnedKeys []string
	// PinnedPrefixes is a list of key prefixes whose values must not drift.
	PinnedPrefixes []string
}

// PinViolation describes a drift result that violates a pin constraint.
type PinViolation struct {
	Key    string
	Status string
	Reason string
}

// Pin inspects results and returns violations for any key that is pinned
// but has drifted (i.e. has a status other than "match").
func Pin(results []Result, opts PinOptions) []PinViolation {
	if len(results) == 0 {
		return nil
	}

	pinnedSet := make(map[string]struct{}, len(opts.PinnedKeys))
	for _, k := range opts.PinnedKeys {
		pinnedSet[strings.ToLower(k)] = struct{}{}
	}

	var violations []PinViolation
	for _, r := range results {
		if strings.EqualFold(r.Status, "match") {
			continue
		}
		if isPinned(r.Key, pinnedSet, opts.PinnedPrefixes) {
			violations = append(violations, PinViolation{
				Key:    r.Key,
				Status: r.Status,
				Reason: pinReason(r.Key, pinnedSet, opts.PinnedPrefixes),
			})
		}
	}

	sort.Slice(violations, func(i, j int) bool {
		return violations[i].Key < violations[j].Key
	})
	return violations
}

func isPinned(key string, exact map[string]struct{}, prefixes []string) bool {
	if _, ok := exact[strings.ToLower(key)]; ok {
		return true
	}
	for _, p := range prefixes {
		if strings.HasPrefix(strings.ToLower(key), strings.ToLower(p)) {
			return true
		}
	}
	return false
}

func pinReason(key string, exact map[string]struct{}, prefixes []string) string {
	if _, ok := exact[strings.ToLower(key)]; ok {
		return "exact key pin"
	}
	for _, p := range prefixes {
		if strings.HasPrefix(strings.ToLower(key), strings.ToLower(p)) {
			return "prefix pin: " + p
		}
	}
	return ""
}
