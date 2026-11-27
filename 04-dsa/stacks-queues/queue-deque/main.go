// Topic: queue-deque | Concept: stacks-queues | Stage: 04-dsa
// Run:   go run ./04-dsa/stacks-queues/queue-deque
//
// Explore:
// * Queue = FIFO: enqueue at the back, dequeue from the front
// * q = q[1:] looks O(1) but keeps the old array alive; a ring buffer reuses space
// * Ring buffer: head index + count over a fixed array, wrap with % cap, grow by doubling
// * Deque (double-ended queue): push/pop at both ends in O(1) amortized
// * Uses: BFS, task queues, sliding-window maximum, undo/redo
// * Buffered channels are Go's concurrency-safe queue between goroutines
package main

import (
	"fmt"
	"os"
	"slices"
)

// Deque is a growable ring buffer. Queue operations are PushBack + PopFront.
type Deque[T any] struct {
	buf  []T
	head int // index of the first element
	n    int // number of elements
}

func (d *Deque[T]) Len() int { return d.n }

func (d *Deque[T]) grow() {
	newCap := max(4, 2*len(d.buf))
	nb := make([]T, newCap)
	for i := range d.n {
		nb[i] = d.buf[(d.head+i)%len(d.buf)]
	}
	d.buf, d.head = nb, 0
}

func (d *Deque[T]) PushBack(v T) {
	if d.n == len(d.buf) {
		d.grow()
	}
	d.buf[(d.head+d.n)%len(d.buf)] = v
	d.n++
}

func (d *Deque[T]) PushFront(v T) {
	if d.n == len(d.buf) {
		d.grow()
	}
	d.head = (d.head - 1 + len(d.buf)) % len(d.buf)
	d.buf[d.head] = v
	d.n++
}

func (d *Deque[T]) PopFront() (T, bool) {
	var zero T
	if d.n == 0 {
		return zero, false
	}
	v := d.buf[d.head]
	d.buf[d.head] = zero
	d.head = (d.head + 1) % len(d.buf)
	d.n--
	return v, true
}

func (d *Deque[T]) PopBack() (T, bool) {
	var zero T
	if d.n == 0 {
		return zero, false
	}
	i := (d.head + d.n - 1) % len(d.buf)
	v := d.buf[i]
	d.buf[i] = zero
	d.n--
	return v, true
}

func (d *Deque[T]) Values() []T {
	out := make([]T, d.n)
	for i := range d.n {
		out[i] = d.buf[(d.head+i)%len(d.buf)]
	}
	return out
}

func main() {
	fmt.Println("Topic: queue and deque (ring buffer)")
	failed := false
	expect := func(label string, got, want []int) {
		status := "ok"
		if !slices.Equal(got, want) {
			status = fmt.Sprintf("MISMATCH (want %v)", want)
			failed = true
		}
		fmt.Printf("  %-28s %v %s\n", label, got, status)
	}

	// As a FIFO queue
	var q Deque[int]
	for i := 1; i <= 5; i++ {
		q.PushBack(i)
	}
	first, _ := q.PopFront()
	second, _ := q.PopFront()
	fmt.Println("queue: dequeued", first, "then", second)
	expect("queue after 2 dequeues", q.Values(), []int{3, 4, 5})

	// Wrap-around: the ring reuses freed slots before growing
	q.PushBack(6)
	q.PushBack(7)
	expect("after wrap-around pushes", q.Values(), []int{3, 4, 5, 6, 7})
	fmt.Println("  backing array cap:", len(q.buf), "| head index:", q.head)

	// As a deque
	var d Deque[int]
	d.PushBack(2)
	d.PushFront(1)
	d.PushBack(3)
	d.PushFront(0)
	expect("deque front/back pushes", d.Values(), []int{0, 1, 2, 3})
	back, _ := d.PopBack()
	front, _ := d.PopFront()
	fmt.Println("  popped back", back, "and front", front)
	expect("deque after pops", d.Values(), []int{1, 2})

	// Compare against a plain slice model across many mixed operations
	var big Deque[int]
	var model []int
	for i := range 1000 {
		switch i % 4 {
		case 0, 1:
			big.PushBack(i)
			model = append(model, i)
		case 2:
			big.PushFront(i)
			model = append([]int{i}, model...)
		case 3:
			big.PopFront()
			model = model[1:]
		}
	}
	ok := slices.Equal(big.Values(), model)
	fmt.Printf("  1000 mixed operations match a slice model: %v (len %d)\n", ok, big.Len())
	if !ok {
		failed = true
	}

	// Channel as a queue between goroutines
	ch := make(chan string, 3)
	ch <- "job-1"
	ch <- "job-2"
	close(ch)
	for job := range ch {
		fmt.Println("  channel queue got:", job)
	}

	if failed {
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}
