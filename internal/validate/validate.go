// Package validate checks KV pairs against user-defined rules and reports
// violations as structured results.
package validate

import (
	"strings"

	"github.com/hashicorp/consul/api"
)

// Severity indicates how serious a violation is.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
)

// Rule describes a single validation constraint.
type Rule struct {
	// Prefix restricts the rule to keys that start with this value.
	Prefix string `yaml:"prefix"`
	// MaxLength is the maximum allowed byte length of the value (0 = unlimited).
	MaxLength int `yaml:"max_length"`
	// Required causes a violation when the value is empty.
	Required bool `yaml:"required"`
	// Severity is the level assigned to violations from this rule.
	Severity Severity `yaml:"severity"`
}

// Violation is a single rule breach for a KV pair.
type Violation struct {
	Key      string
	Rule     Rule
	Message  string
	Severity Severity
}

// Apply evaluates pairs against rules and returns all violations found.
func Apply(pairs []*api.KVPair, rules []Rule) []Violation {
	var out []Violation
	for _, p := range pairs {
		if p == nil {
			continue
		}
		for _, r := range rules {
			if r.Prefix != "" && !strings.HasPrefix(p.Key, r.Prefix) {
				continue
			}
			sev := r.Severity
			if sev == "" {
				sev = SeverityWarning
			}
			if r.Required && len(p.Value) == 0 {
				out = append(out, Violation{
					Key:      p.			Rule:     r,
					Message:  "value is required but empty",
					Severity: sev,
				})
			}
			if r.MaxLength > 0 && len(p.Value) > r.MaxLength {
				out = append(out, Violation{
					Key:      p.Key,
					Rule:     r,
					Message:  "value exceeds maximum length",
					Severity: sev,
				})
			}
		}
	}
	return out
}
