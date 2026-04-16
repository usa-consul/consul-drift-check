// Package throttle provides rate-limiting for Consul API requests.
package throttle

import (
	"context"
	"time"
)

// Throttle controls the rate of operations using a token bucket approach.
type Throttle struct {
	ticker  *time.Ticker
	tokens  chan struct{}
	done    chan struct{}
}

// New creates a Throttle that allows up to rps operations per second.
func New(rps int) *Throttle {
	if rps <= 0 {
		rps = 1
	}
	interval := time.Second / time.Duration(rps)
	t := &Throttle{
		ticker: time.NewTicker(interval),
		tokens: make(chan struct{}, rps),
		done:   make(chan struct{}),
	}
	go t.fill()
	return t
}

func (t *Throttle) fill() {
	for {
		select {
		case <-t.ticker.C:
			select {
			case t.tokens <- struct{}{}:
			default:
			}
		case <-t.done:
			return
		}
	}
}

// Wait blocks until a token is available or ctx is cancelled.
func (t *Throttle) Wait(ctx context.Context) error {
	select {
	case <-t.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Stop releases resources held by the Throttle.
func (t *Throttle) Stop() {
	t.ticker.Stop()
	close(t.done)
}
