// Package sample records point-in-time snapshots of Consul KV prefix data.
//
// Entries are appended as newline-delimited JSON to a local file, allowing
// drift trends to be analysed over time. Use Since to filter entries by
// a start time when building rolling-window comparisons.
package sample
