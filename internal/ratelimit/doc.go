// Package ratelimit provides per-datacenter token-bucket rate limiting for
// outbound Consul API requests made by consul-drift-check.
//
// Each datacenter gets an independent bucket that refills at the configured
// requests-per-second rate. Callers block on Wait until a token is available
// or the supplied context is cancelled.
//
// Usage:
//
//	limiter := ratelimit.New(10) // 10 RPS per datacenter
//	defer limiter.Stop()
//
//	if err := limiter.Wait(ctx, "dc1"); err != nil {
//		return err
//	}
//	// proceed with Consul API call
package ratelimit
