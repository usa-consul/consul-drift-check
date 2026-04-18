// Package annotate provides rule-based annotation of diff results.
//
// Rules are evaluated in order; the first prefix match wins. Annotations are
// arbitrary key/value pairs that downstream consumers (reports, exports, etc.)
// may use for display or filtering purposes.
package annotate
