// Topic: buffered-unbuffered | Concept: channels | Stage: 03-advanced
// Run:   go run ./03-advanced/channels/buffered-unbuffered
//
// Explore:
// * "Don't communicate by sharing memory; share memory by communicating"
// * Unbuffered channel: send blocks until a receiver takes the value
// * Buffered channel: send blocks only when the buffer is full
// * close(ch) signals "no more values"; range over a channel stops at close
// * v, ok := <-ch: ok is false once closed and drained
// * Direction types: chan<- int (send only), <-chan int (receive only)
// * Sending on a closed channel panics; a receive-forever with no sender deadlocks
package main

import "fmt"

// producer can only send; the caller can only receive
func producer(n int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 1; i <= n; i++ {
			out <- i
		}
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range in {
			out <- v * v
		}
	}()
	return out
}

func main() {
	fmt.Println("Topic: buffered and unbuffered channels")

	// Unbuffered: a handshake between two goroutines
	done := make(chan string)
	go func() {
		done <- "work finished"
	}()
	fmt.Println("unbuffered received:", <-done)

	// Buffered: up to cap values without a receiver
	jobs := make(chan string, 3)
	jobs <- "a"
	jobs <- "b"
	fmt.Printf("buffered: len=%d cap=%d\n", len(jobs), cap(jobs))
	jobs <- "c"
	// jobs <- "d" // would block forever here: fatal error: all goroutines are asleep - deadlock!
	close(jobs)
	for j := range jobs {
		fmt.Println("  drained", j)
	}
	v, ok := <-jobs
	fmt.Printf("receive after close: %q, ok=%v\n", v, ok)

	// Pipeline: producer -> square -> main
	fmt.Print("pipeline squares: ")
	for sq := range square(producer(5)) {
		fmt.Print(sq, " ")
	}
	fmt.Println()
}
