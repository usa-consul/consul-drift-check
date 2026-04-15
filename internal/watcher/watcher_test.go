package watcher_test

import (
	"context"
	"errors"
	"log"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/your-org/consul-drift-check/internal/watcher"
)

func silentLogger() *log.Logger {
	return log.New(os.Discard, "", 0)
}

func TestWatcher_InvokesHandlerImmediately(t *testing.T) {
	var count int32
	handler := func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	w := watcher.New(10*time.Second, handler, silentLogger())
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	w.Run(ctx) //nolint:errcheck

	if atomic.LoadInt32(&count) < 1 {
		t.Fatal("expected handler to be called at least once immediately")
	}
}

func TestWatcher_TicksOnInterval(t *testing.T) {
	var count int32
	handler := func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	w := watcher.New(20*time.Millisecond, handler, silentLogger())
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Millisecond)
	defer cancel()

	w.Run(ctx) //nolint:errcheck

	got := atomic.LoadInt32(&count)
	if got < 3 {
		t.Fatalf("expected at least 3 invocations, got %d", got)
	}
}

func TestWatcher_HandlerErrorDoesNotStop(t *testing.T) {
	var count int32
	handler := func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return errors.New("simulated error")
	}

	w := watcher.New(20*time.Millisecond, handler, silentLogger())
	ctx, cancel := context.WithTimeout(context.Background(), 70*time.Millisecond)
	defer cancel()

	w.Run(ctx) //nolint:errcheck

	if atomic.LoadInt32(&count) < 2 {
		t.Fatal("expected watcher to continue after handler error")
	}
}

func TestWatcher_StopsOnContextCancel(t *testing.T) {
	handler := func(ctx context.Context) error { return nil }

	w := watcher.New(5*time.Second, handler, silentLogger())
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() {
		done <- w.Run(ctx)
	}()

	cancel()

	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("watcher did not stop after context cancellation")
	}
}
