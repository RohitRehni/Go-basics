// Topic: once | Concept: sync-patterns | Stage: 03-advanced
// Run:   go run ./03-advanced/sync-patterns/once
// Race:  go run -race ./03-advanced/sync-patterns/once
//
// Explore:
// * sync.Once runs a function exactly once, even from many goroutines
// * Lazy initialisation of expensive things (config, DB pool, regex)
// * sync.OnceValue / sync.OnceValues (Go 1.21+) return the cached result
// * If the once function panics, later calls do NOT retry it
// * sync.Pool for reusing short-lived objects (buffers)
package main

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
)

type Config struct {
	DSN string
}

var (
	loadCount atomic.Int32
	getConfig = sync.OnceValue(func() *Config {
		loadCount.Add(1)
		fmt.Println("  loading config (expensive, happens once)")
		return &Config{DSN: "postgres://localhost/learning"}
	})
)

var bufPool = sync.Pool{
	New: func() any { return new(bytes.Buffer) },
}

func render(name string) string {
	buf := bufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bufPool.Put(buf)
	}()
	buf.WriteString("hello, ")
	buf.WriteString(name)
	return buf.String()
}

func main() {
	fmt.Println("Topic: sync.Once, OnceValue and Pool")

	var once sync.Once
	var wg sync.WaitGroup
	for i := range 5 {
		wg.Go(func() {
			once.Do(func() { fmt.Println("  once.Do ran (from one of 5 goroutines)") })
			_ = i
		})
	}
	wg.Wait()

	for range 10 {
		wg.Go(func() { _ = getConfig() })
	}
	wg.Wait()
	fmt.Println("config DSN:", getConfig().DSN)
	fmt.Println("times loaded:", loadCount.Load())

	fmt.Println("pool render:", render("go"))
	fmt.Println("pool render:", render("rust"))
}
