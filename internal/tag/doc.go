// Package tag provides prefix-based tagging of drift results.
//
// Rules are evaluated in order; the first matching prefix wins.
// Tags can be used downstream by exporters, alerting, or grouping logic
// to attach environment, team, or criticality metadata to each result.
package tag
