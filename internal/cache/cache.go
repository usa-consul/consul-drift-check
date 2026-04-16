// Package cache provides a simple in-memory cache for Consul KV entries
// to reduce redundant API calls during drift checks.
package cache

import (
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
)

// Entry holds cached KV pairs along with metadata.
type Entry struct {
	Pairs  api.KVPairs
	CachedAt time.Time
}

// Cache is a thread-safe in-memory store for KV list results.
type Cache struct {
	mu      sync.RWMutex
	store   map[string]Entry
	ttl     time.Duration
}

// New creates a Cache with the given TTL duration.
func New(ttl time.Duration) *Cache {
	return &Cache{
		store: make(map[string]Entry),
		ttl:   ttl,
	}
}

// Get returns cached KV pairs for the given key, and whether the entry
// was found and is still valid.
func (c *Cache) Get(key string) (api.KVPairs, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.store[key]
	if !ok {
		return nil, false
	}
	if time.Since(entry.CachedAt) > c.ttl {
		return nil, false
	}
	return entry.Pairs, true
}

// Set stores KV pairs under the given key.
func (c *Cache) Set(key string, pairs api.KVPairs) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store[key] = Entry{
		Pairs:    pairs,
		CachedAt: time.Now(),
	}
}

// Invalidate removes the entry for the given key.
func (c *Cache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.store, key)
}

// Flush removes all entries from the cache.
func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store = make(map[string]Entry)
}
