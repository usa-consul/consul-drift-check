package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"github.com/your-org/consul-drift-check/internal/ratelimit"
)

func TestWait_AcquiresToken(t *testing.T) {
	l := ratelimit.New(10)
	defer l.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := l.Wait(ctx, "dc1"); err != nil {
		t.Fatalf("expected token, got error: %v", err)
	}
}

func TestWait_CancelledContext(t *testing.T) {
	l := ratelimit.New(1)
	defer l.Stop()

	// Drain the first token so the bucket is empty.
	ctx := context.Background()
	_ = l.Wait(ctx, "dc-slow")

	// Now cancel immediately — should get context error.
	ctx2, cancel := context.WithCancel(context.Background())
	cancel()

	err := l.Wait(ctx2, "dc-slow")
	if err == nil {
		t.Fatal("expected error for cancelled context, got nil")
	}
}

func TestNew_ZeroRPS_DefaultsToOne(t *testing.T) {
	l := ratelimit.New(0)
	defer l.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := l.Wait(ctx, "dc1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWait_MultipleDCsAreIndependent(t *testing.T) {
	l := ratelimit.New(5)
	defer l.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	for _, dc := range []string{"dc1", "dc2", "dc3"} {
		if err := l.Wait(ctx, dc); err != nil {
			t.Fatalf("dc %s: unexpected error: %v", dc, err)
		}
	}
}

func TestStop_DoesNotPanic(t *testing.T) {
	l := ratelimit.New(5)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_ = l.Wait(ctx, "dc1")
	l.Stop()
	l.Stop() // second call must not panic
}
