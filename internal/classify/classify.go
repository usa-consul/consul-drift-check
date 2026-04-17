// Package classify assigns a severity level to each drift result
// based on configurable key patterns and change types.
package classify

import (
	"strings"

	"github.com/your-org/consul-drift-check/internal/diff"
)

// Level represents the severity of a drift result.
type Level string

const (
	LevelInfo     Level = "info"
	LevelWarning  Level = "warning"
	LevelCritical Level = "critical"
)

// Rule maps a key prefix to a severity level.
type Rule struct {
	Prefix string
	Level  Level
}

// Result pairs a diff result with its assigned severity.
type Result struct {
	diff.Result
	Level Level
}

// Apply classifies each diff result using the provided rules.
// Rules are evaluated in order; the first matching prefix wins.
// Results that match no rule default to LevelInfo.
func Apply(results []diff.Result, rules []Rule) []Result {
	out := make([]Result, 0, len(results))
	for _, r := range results {
		out = append(out, Result{
			Result: r,
			Level:  resolve(r.Key, rules),
		})
	}
	return out
}

func resolve(key string, rules []Rule) Level {
	for _, rule := range rules {
		if strings.HasPrefix(key, rule.Prefix) {
			return rule.Level
		}
	}
	return LevelInfo
}
