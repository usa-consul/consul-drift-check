// Package ratelimit provides per-datacenter rate limiting for Consul API calls.
package ratelimit

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Limiter enforces a maximum number of requests per second per datacenter.
type Limiter struct {
	mu      sync.Mutex
	buckets map[string]chan struct{}
	rps     int
	stopChs map[string]chan struct{}
}

// New creates a Limiter allowing rps requests per second per datacenter.
func New(rps int) *Limiter {
	if rps <= 0 {
		rps = 1
	}
	return &Limiter{
		buckets: make(map[string]chan struct{}),
		stopChs: make(map[string]chan struct{}),
		rps:     rps,
	}
}

// Wait blocks until a token is available for the given datacenter or ctx is done.
func (l *Limiter) Wait(ctx context.Context, dc string) error {
	l.mu.Lock()
	if _, ok := l.buckets[dc]; !ok {
		l.start(dc)
	}
	bucket := l.buckets[dc]
	l.mu.Unlock()

	select {
	case <-bucket:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("ratelimit: context cancelled waiting for dc %q: %w", dc, ctx.Err())
	}
}

// start initialises the token bucket goroutine for a datacenter.
// Must be called with l.mu held.
func (l *Limiter) start(dc string) {
	bucket := make(chan struct{}, l.rps)
	stop := make(chan struct{})
	l.buckets[dc] = bucket
	l.stopChs[dc] = stop

	interval := time.Second / time.Duration(l.rps)
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				select {
				case bucket <- struct{}{}:
				default:
				}
			case <-stop:
				return
			}
		}
	}()
}

// StopDatacenter releases resources for a single datacenter.
func (l *Limiter) StopDatacenter(dc string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if ch, ok := l.stopChs[dc]; ok {
		close(ch)
		delete(l.stopChs, dc)
		delete(l.buckets, dc)
	}
}

// Stop releases resources for all datacenters.
func (l *Limiter) Stop() {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, ch := range l.stopChs {
		close(ch)
	}
	l.stopChs = make(map[string]chan struct{})
	l.buckets = make(map[string]chan struct{})
}
