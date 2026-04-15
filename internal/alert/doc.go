// Package alert evaluates consul-drift-check results against configurable
// thresholds and produces OK / WARNING / CRITICAL severity levels.
//
// Usage:
//
//	summary := metrics.Collect(results)
//	thresholds := alert.Thresholds{Warning: 5, Critical: 20}
//	result := alert.Evaluate(summary, thresholds)
//	alert.Write(os.Stdout, result)
//
// The alert level is determined by the total number of drifted keys
// (added + removed + modified) relative to the configured thresholds.
package alert
