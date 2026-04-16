// Package retry provides simple retry logic with backoff for transient errors.
package retry

import (
	"context"
	"errors"
	"time"
)

// Policy defines retry behaviour.
type Policy struct {
	MaxAttempts int
	InitialDelay time.Duration
	Multiplier float64
}

// DefaultPolicy returns a sensible default retry policy.
func DefaultPolicy() Policy {
	return Policy{
		MaxAttempts:  3,
		InitialDelay: 200 * time.Millisecond,
		Multiplier:   2.0,
	}
}

// ErrMaxAttemptsReached is returned when all attempts are exhausted.
var ErrMaxAttemptsReached = errors.New("retry: max attempts reached")

// Do executes fn according to p, retrying on non-nil errors.
// The context is checked before each attempt.
func Do(ctx context.Context, p Policy, fn func() error) error {
	if p.MaxAttempts < 1 {
		p.MaxAttempts = 1
	}
	if p.Multiplier <= 0 {
		p.Multiplier = 1
	}

	delay := p.InitialDelay
	var last error

	for attempt := 0; attempt < p.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		if last = fn(); last == nil {
			return nil
		}

		if attempt < p.MaxAttempts-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
			delay = time.Duration(float64(delay) * p.Multiplier)
		}
	}

	return ErrMaxAttemptsReached
}
