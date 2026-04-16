// Package redact provides utilities to mask sensitive KV values before
// they are written to reports, audit logs, or transmitted over the wire.
//
// Matching is performed case-insensitively against the full key path or
// the final path segment, so registering "password" will redact both
// "password" and "db/password".
package redact
