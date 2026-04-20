package diff

import "fmt"

// QuotaOptions configures drift quota thresholds per prefix.
type QuotaOptions struct {
	// Rules maps a key prefix to a maximum allowed drift count.
	Rules map[string]int
	// DefaultMax is used when no rule matches. Zero means unlimited.
	DefaultMax int
}

// QuotaViolation describes a prefix that has exceeded its drift quota.
type QuotaViolation struct {
	Prefix  string
	Allowed int
	Actual  int
}

// String returns a human-readable description of the violation.
func (v QuotaViolation) String() string {
	return fmt.Sprintf("prefix %q: allowed %d, actual %d", v.Prefix, v.Allowed, v.Actual)
}

// EvaluateQuota checks whether drift counts per top-level prefix exceed
// configured quotas. It returns a slice of violations (nil if none).
func EvaluateQuota(results []Result, opts QuotaOptions) []QuotaViolation {
	if len(results) == 0 {
		return nil
	}

	counts := make(map[string]int)
	for _, r := range results {
		seg := topQuotaSegment(r.Key)
		counts[seg]++
	}

	var violations []QuotaViolation
	for prefix, actual := range counts {
		max := opts.DefaultMax
		if v, ok := opts.Rules[prefix]; ok {
			max = v
		}
		if max > 0 && actual > max {
			violations = append(violations, QuotaViolation{
				Prefix:  prefix,
				Allowed: max,
				Actual:  actual,
			})
		}
	}

	sortQuotaViolations(violations)
	return violations
}

func topQuotaSegment(key string) string {
	if len(key) > 0 && key[0] == '/' {
		key = key[1:]
	}
	for i, c := range key {
		if c == '/' {
			return key[:i]
		}
	}
	return key
}

func sortQuotaViolations(vs []QuotaViolation) {
	for i := 1; i < len(vs); i++ {
		for j := i; j > 0 && vs[j].Prefix < vs[j-1].Prefix; j-- {
			vs[j], vs[j-1] = vs[j-1], vs[j]
		}
	}
}
