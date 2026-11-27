// Topic: waitgroup | Concept: goroutines | Stage: 03-advanced
// Run:   go run ./03-advanced/goroutines/waitgroup
// Race:  go run -race ./03-advanced/goroutines/waitgroup   (race detector)
//
// Explore:
// * `go f()` starts f in a new goroutine: cheap, managed by the Go runtime
// * main does not wait for goroutines: when main returns, they all stop
// * sync.WaitGroup: Add before starting, Done when finished, Wait to block
// * wg.Go(func) (Go 1.25+) does Add/Done for you
// * Goroutines run in any order: never rely on print order
// * Collect results by index (each goroutine owns its slot) to avoid races
package main

import (
	"fmt"
	"sync"
	"time"
)

func fetch(id int) string {
	time.Sleep(time.Duration(30-id*5) * time.Millisecond) // later ids finish first
	return fmt.Sprintf("result-%d", id)
}

func main() {
	fmt.Println("Topic: goroutines and sync.WaitGroup")

	start := time.Now()
	results := make([]string, 5)

	var wg sync.WaitGroup
	for i := range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results[i] = fetch(i) // each goroutine writes only its own index
		}()
	}
	wg.Wait()

	fmt.Println("results in index order:", results)
	fmt.Printf("5 fetches of up to 30ms took about %dms (they ran concurrently)\n",
		time.Since(start).Round(10*time.Millisecond).Milliseconds())

	// wg.Go: the Go 1.25+ shortcut
	var wg2 sync.WaitGroup
	squares := make([]int, 4)
	for i := range squares {
		wg2.Go(func() { squares[i] = i * i })
	}
	wg2.Wait()
	fmt.Println("squares via wg.Go:", squares)

	// Without Wait, main could exit before this prints:
	// go fmt.Println("may never appear")
}
