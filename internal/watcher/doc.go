// Package watcher provides a periodic scheduler for running drift checks
// at a configurable interval. It is designed to integrate with the
// consul-drift-check pipeline, accepting a Handler function that encapsulates
// a full drift check cycle.
//
// Usage:
//
//	w := watcher.New(5*time.Minute, myHandler, nil)
//	if err := w.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
//		log.Fatal(err)
//	}
//
// The handler is called immediately on startup and then on every tick.
// Handler errors are logged but do not stop the watcher.
package watcher
