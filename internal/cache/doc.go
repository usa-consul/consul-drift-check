// Package cache provides a lightweight in-memory cache for Consul KV results.
//
// It is intended to reduce redundant API calls when the same KV prefix is
// queried multiple times within a single drift-check run (e.g. in watch mode).
//
// Each cache entry has a configurable TTL; stale entries are treated as
// misses and the caller is responsible for refreshing them.
//
// Example usage:
//
//	c := cache.New(30 * time.Second)
//	if pairs, ok := c.Get(prefix); ok {
//		return pairs, nil
//	}
//	pairs, err := kvClient.ListPrefix(ctx, prefix)
//	if err == nil {
//		c.Set(prefix, pairs)
//	}
package cache
