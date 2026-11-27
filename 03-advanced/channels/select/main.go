// Topic: select | Concept: channels | Stage: 03-advanced
// Run:   go run ./03-advanced/channels/select
//
// Explore:
// * select waits on several channel operations, runs the first one ready
// * If several are ready, one is picked at random
// * time.After for a timeout case
// * default makes select non-blocking (try-send / try-receive)
// * A nil channel blocks forever: set a channel to nil to disable its case
// * for + select is the usual event loop
package main

import (
	"fmt"
	"time"
)

func after(d time.Duration, msg string) <-chan string {
	ch := make(chan string, 1)
	go func() {
		time.Sleep(d)
		ch <- msg
	}()
	return ch
}

func main() {
	fmt.Println("Topic: select")

	// First one wins
	select {
	case m := <-after(20*time.Millisecond, "fast service"):
		fmt.Println("winner:", m)
	case m := <-after(80*time.Millisecond, "slow service"):
		fmt.Println("winner:", m)
	}

	// Timeout
	select {
	case m := <-after(100*time.Millisecond, "too slow"):
		fmt.Println("got:", m)
	case <-time.After(30 * time.Millisecond):
		fmt.Println("timeout after 30ms")
	}

	// Non-blocking with default
	queue := make(chan int, 1)
	for i := 1; i <= 2; i++ {
		select {
		case queue <- i:
			fmt.Println("queued", i)
		default:
			fmt.Println("queue full, dropped", i)
		}
	}

	// Event loop: merge two channels until both are closed
	a := make(chan string)
	b := make(chan string)
	go func() { defer close(a); a <- "a1"; a <- "a2" }()
	go func() { defer close(b); b <- "b1" }()

	received := 0
	for a != nil || b != nil {
		select {
		case v, ok := <-a:
			if !ok {
				a = nil // disable this case
				continue
			}
			received++
			fmt.Println("from a:", v)
		case v, ok := <-b:
			if !ok {
				b = nil
				continue
			}
			received++
			fmt.Println("from b:", v)
		}
	}
	fmt.Println("total received:", received)
}
