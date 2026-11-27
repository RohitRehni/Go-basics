// Topic: atomic | Concept: sync-patterns | Stage: 03-advanced
// Run:   go run ./03-advanced/sync-patterns/atomic
// Race:  go run -race ./03-advanced/sync-patterns/atomic
//
// Explore:
// * sync/atomic: lock-free operations on single values
// * Typed atomics: atomic.Int64, atomic.Bool, atomic.Pointer[T]
// * Add, Load, Store, Swap, CompareAndSwap
// * Good for counters and flags; use a Mutex when several fields change together
// * atomic.Value / Pointer for swapping a whole config at runtime
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Settings struct {
	LogLevel string
}

func main() {
	fmt.Println("Topic: sync/atomic")

	var requests atomic.Int64
	var wg sync.WaitGroup
	for range 1000 {
		wg.Go(func() { requests.Add(1) })
	}
	wg.Wait()
	fmt.Println("requests after 1000 goroutines:", requests.Load())

	var ready atomic.Bool
	fmt.Println("ready before:", ready.Load())
	ready.Store(true)
	fmt.Println("ready after Store(true):", ready.Load())

	var state atomic.Int32 // 0 = idle, 1 = running
	fmt.Println("CAS idle->running:", state.CompareAndSwap(0, 1))
	fmt.Println("CAS idle->running again:", state.CompareAndSwap(0, 1), "(already running)")
	old := state.Swap(0)
	fmt.Println("Swap back to idle, old value:", old)

	// Only one of many goroutines wins the CAS
	var winner atomic.Int32
	var wins atomic.Int32
	for i := 1; i <= 10; i++ {
		wg.Go(func() {
			if winner.CompareAndSwap(0, int32(i)) {
				wins.Add(1)
			}
		})
	}
	wg.Wait()
	fmt.Println("goroutines that won the CAS:", wins.Load())

	var cfg atomic.Pointer[Settings]
	cfg.Store(&Settings{LogLevel: "info"})
	cfg.Store(&Settings{LogLevel: "debug"}) // hot reload
	fmt.Println("current log level:", cfg.Load().LogLevel)
}
