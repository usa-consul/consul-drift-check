package throttle_test

import (
	"context"
	"testing"
	"time"

	"github.com/your-org/consul-drift-check/internal/throttle"
)

func TestWait_AcquiresToken(t *testing.T) {
	th := throttle.New(100)
	defer th.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := th.Wait(ctx); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestWait_CancelledContext(t *testing.T) {
	th := throttle.New(1)
	defer th.Stop()

	// Drain the first token.
	ctx := context.Background()
	_ = th.Wait(ctx)

	// Now cancel immediately before a new token is available.
	cancel_ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := th.Wait(cancel_ctx)
	if err == nil {
		t.Fatal("expected context cancellation error")
	}
}

func TestNew_ZeroRPS_DefaultsToOne(t *testing.T) {
	th := throttle.New(0)
	defer th.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := th.Wait(ctx); err != nil {
		t.Fatalf("expected token within timeout, got %v", err)
	}
}

func TestStop_DoesNotPanic(t *testing.T) {
	th := throttle.New(10)
	th.Stop()
}
