// Topic: mutex | Concept: goroutines | Stage: 03-advanced
// Run:   go run ./03-advanced/goroutines/mutex
// Race:  go run -race ./03-advanced/goroutines/mutex
//
// Explore:
// * A data race: two goroutines touch the same variable, at least one writes
// * sync.Mutex: Lock / Unlock around the shared data (defer Unlock is safest)
// * Put the mutex next to the data it guards, inside a struct
// * sync.RWMutex: many readers OR one writer
// * Never copy a struct that contains a mutex (use a pointer)
// * `go run -race` finds races; uncomment the unsafe block to see it
package main

import (
	"fmt"
	"sync"
)

type SafeCounter struct {
	mu     sync.Mutex
	counts map[string]int
}

func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counts[key]++
}

func (c *SafeCounter) Get(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.counts[key]
}

type Cache struct {
	mu   sync.RWMutex
	data map[string]string
}

func (c *Cache) Read(k string) (string, bool) {
	c.mu.RLock() // many readers at once
	defer c.mu.RUnlock()
	v, ok := c.data[k]
	return v, ok
}

func (c *Cache) Write(k, v string) {
	c.mu.Lock() // exclusive
	defer c.mu.Unlock()
	c.data[k] = v
}

func main() {
	fmt.Println("Topic: sync.Mutex and sync.RWMutex")

	counter := SafeCounter{counts: map[string]int{}}
	var wg sync.WaitGroup
	for range 100 {
		wg.Go(func() { counter.Inc("hits") })
	}
	wg.Wait()
	fmt.Println("hits after 100 goroutines (always 100):", counter.Get("hits"))

	cache := &Cache{data: map[string]string{}}
	cache.Write("lang", "go")
	for range 5 {
		wg.Go(func() { cache.Read("lang") })
	}
	wg.Wait()
	v, _ := cache.Read("lang")
	fmt.Println("cache lang =", v)

	// UNSAFE: a data race. Uncomment and run with -race to see the report.
	// total := 0
	// for range 100 {
	// 	wg.Go(func() { total++ })
	// }
	// wg.Wait()
	// fmt.Println("racy total (often < 100):", total)
}
