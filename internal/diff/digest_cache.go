package diff

import (
	"sync"
	"time"
)

// digestCacheEntry stores a DigestEntry alongside its expiry.
type digestCacheEntry struct {
	value     DigestEntry
	expiresAt time.Time
}

// DigestCache is a thread-safe in-memory store for DigestEntry values.
type DigestCache struct {
	mu  sync.RWMutex
	ttl time.Duration
	m   map[string]digestCacheEntry
}

// NewDigestCache creates a DigestCache with the given TTL.
func NewDigestCache(ttl time.Duration) *DigestCache {
	if ttl <= 0 {
		ttl = 30 * time.Second
	}
	return &DigestCache{ttl: ttl, m: make(map[string]digestCacheEntry)}
}

// Set stores a DigestEntry under the given key.
func (c *DigestCache) Set(key string, e DigestEntry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key] = digestCacheEntry{value: e, expiresAt: time.Now().Add(c.ttl)}
}

// Get retrieves a DigestEntry. The second return value is false on a miss or
// when the entry has expired.
func (c *DigestCache) Get(key string) (DigestEntry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.m[key]
	if !ok || time.Now().After(e.expiresAt) {
		return DigestEntry{}, false
	}
	return e.value, true
}

// Invalidate removes a single entry from the cache.
func (c *DigestCache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.m, key)
}
