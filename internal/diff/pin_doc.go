// Package diff provides utilities for comparing Consul KV namespaces.
//
// # Pin
//
// Pin enforces that specific keys or key prefixes must not drift between
// the source and destination datacenters. Any result whose key matches a
// pinned constraint and whose status is not "match" is reported as a
// PinViolation.
//
// Pinning is useful for critical configuration keys — such as TLS
// certificates, bootstrap tokens, or feature flags — where any divergence
// should be treated as a hard failure regardless of other drift thresholds.
//
// Example:
//
//	violations := diff.Pin(results, diff.PinOptions{
//		PinnedKeys:     []string{"config/tls/cert"},
//		PinnedPrefixes: []string{"config/bootstrap/"},
//	})
package diff
