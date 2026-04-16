package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/your-org/consul-drift-check/internal/retry"
)

var errTransient = errors.New("transient error")

func TestDo_SucceedsOnFirstAttempt(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), retry.DefaultPolicy(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesAndSucceeds(t *testing.T) {
	calls := 0
	p := retry.Policy{MaxAttempts: 3, InitialDelay: time.Millisecond, Multiplier: 1}
	err := retry.Do(context.Background(), p, func() error {
		calls++
		if calls < 3 {
			return errTransient
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil after retries, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ExhaustsAttempts(t *testing.T) {
	p := retry.Policy{MaxAttempts: 2, InitialDelay: time.Millisecond, Multiplier: 1}
	err := retry.Do(context.Background(), p, func() error {
		return errTransient
	})
	if !errors.Is(err, retry.ErrMaxAttemptsReached) {
		t.Fatalf("expected ErrMaxAttemptsReached, got %v", err)
	}
}

func TestDo_CancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := retry.Do(ctx, retry.DefaultPolicy(), func() error {
		return errTransient
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestDo_ZeroMaxAttempts_DefaultsToOne(t *testing.T) {
	calls := 0
	p := retry.Policy{MaxAttempts: 0, InitialDelay: time.Millisecond, Multiplier: 1}
	retry.Do(context.Background(), p, func() error { //nolint
		calls++
		return errTransient
	})
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}
