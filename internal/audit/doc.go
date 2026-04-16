// Package audit provides append-only audit logging for consul-drift-check runs.
//
// Each invocation of the drift check can record an Entry containing the
// timestamp, config path, KV prefix checked, drift metrics summary, and any
// error encountered. Entries are stored as newline-delimited JSON (JSONL) so
// the log file can be streamed or tailed without loading it entirely into
// memory.
//
// Usage:
//
//	logger := audit.NewLogger("/var/log/consul-drift-check/audit.jsonl")
//	err := logger.Record(audit.Entry{
//		ConfigPath: "/etc/consul-drift/config.yaml",
//		Prefix:     "services/",
//		Summary:    summary,
//	})
package audit
