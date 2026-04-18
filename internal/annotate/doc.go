// Package annotate provides rule-based annotation of diff results.
//
// Rules are evaluated in order; the first prefix match wins. Annotations are
// arbitrary key/value pairs that downstream consumers (reports, exports, etc.)
// may use for display or filtering purposes.
//
// Basic usage:
//
//	ruleset := annotate.NewRuleSet()
//	ruleset.Add(annotate.Rule{
//		Prefix: "aws/us-east-1",
//		Annotations: map[string]string{"region": "us-east-1", "cloud": "aws"},
//	})
//	annotations := ruleset.Match("aws/us-east-1/service/web")
package annotate
