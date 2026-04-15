// Package watcher provides periodic drift check scheduling.
package watcher

import (
	"context"
	"log"
	"time"
)

// Handler is a function called on each tick to perform a drift check.
type Handler func(ctx context.Context) error

// Watcher runs a Handler on a configurable interval.
type Watcher struct {
	interval time.Duration
	handler  Handler
	logger   *log.Logger
}

// New creates a Watcher with the given interval and handler.
func New(interval time.Duration, handler Handler, logger *log.Logger) *Watcher {
	if logger == nil {
		logger = log.Default()
	}
	return &Watcher{
		interval: interval,
		handler:  handler,
		logger:   logger,
	}
}

// Run starts the watcher loop. It blocks until ctx is cancelled.
// The handler is invoked immediately on start, then on each tick.
func (w *Watcher) Run(ctx context.Context) error {
	if err := w.tick(ctx); err != nil {
		w.logger.Printf("watcher: handler error: %v", err)
	}

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Println("watcher: context cancelled, stopping")
			return ctx.Err()
		case <-ticker.C:
			if err := w.tick(ctx); err != nil {
				w.logger.Printf("watcher: handler error: %v", err)
			}
		}
	}
}

func (w *Watcher) tick(ctx context.Context) error {
	w.logger.Println("watcher: running drift check")
	return w.handler(ctx)
}
