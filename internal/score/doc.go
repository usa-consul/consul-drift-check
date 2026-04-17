// Package score computes a weighted numeric severity score that summarises
// the magnitude of configuration drift detected between two Consul KV
// namespaces.
//
// A higher score indicates more significant drift. Modified keys carry a
// higher weight than keys that are present in only one side, reflecting the
// fact that silent value changes are generally more dangerous than missing
// keys that are easy to spot.
package score
