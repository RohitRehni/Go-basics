// Topic: cancel | Concept: context | Stage: 03-advanced
// Run:   go run ./03-advanced/context/cancel
//
// Explore:
// * context.WithCancel: you decide when to stop
// * Cancelling a parent cancels every child context derived from it
// * One cancel stops many goroutines (fan-out shutdown)
// * context.WithValue for request-scoped data (request ID), not for options
// * context.WithCancelCause / context.Cause to say why it was stopped
package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type ctxKey string

const requestIDKey ctxKey = "requestID"

func poller(ctx context.Context, name string, ticks *int, mu *sync.Mutex) {
	t := time.NewTicker(10 * time.Millisecond)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("  %s stopping (req %v): %v\n", name, ctx.Value(requestIDKey), ctx.Err())
			return
		case <-t.C:
			mu.Lock()
			*ticks++
			mu.Unlock()
		}
	}
}

func main() {
	fmt.Println("Topic: context cancellation")

	parent := context.WithValue(context.Background(), requestIDKey, "req-42")
	ctx, cancel := context.WithCancel(parent)

	var (
		wg    sync.WaitGroup
		mu    sync.Mutex
		ticks int
	)
	for _, name := range []string{"poller-A", "poller-B", "poller-C"} {
		wg.Go(func() { poller(ctx, name, &ticks, &mu) })
	}

	time.Sleep(55 * time.Millisecond)
	fmt.Println("cancelling all pollers with one call")
	cancel()
	wg.Wait()
	fmt.Println("all pollers stopped; ticks recorded > 0:", ticks > 0)

	// Parent cancel reaches children
	root, stopRoot := context.WithCancel(context.Background())
	child, stopChild := context.WithTimeout(root, time.Hour)
	defer stopChild()
	stopRoot()
	<-child.Done()
	fmt.Println("child cancelled by parent:", child.Err())

	// Cancel with a cause
	errShutdown := errors.New("server shutting down")
	cctx, cancelCause := context.WithCancelCause(context.Background())
	cancelCause(errShutdown)
	fmt.Println("ctx.Err():", cctx.Err(), "| context.Cause:", context.Cause(cctx))
}
