// Package diff provides utilities for comparing Consul KV namespaces.
//
// The digest sub-feature computes a deterministic SHA-256 fingerprint over an
// ordered set of KV pairs. Two DigestEntry values with identical Digest strings
// are guaranteed to have the same keys and values, making digest comparison a
// cheap first-pass check before running a full key-by-key diff.
package diff
