// Topic: min-heap | Concept: heaps | Stage: 04-dsa
// Run:   go run ./04-dsa/heaps/min-heap
//
// Explore:
//   - A binary heap is a complete tree stored in a slice:
//     children of i are 2i+1 and 2i+2, parent of i is (i-1)/2
//   - Min-heap rule: every parent <= its children, so the minimum is at index 0
//   - Push: append, then sift up. Pop: move the last item to the root, then sift down. Both O(log n)
//   - Peek O(1); building a heap from n items with heapify is O(n)
//   - Priority queue = heap of (priority, item); top-k largest = min-heap of size k
//   - Heap sort: heapify, then pop n times: O(n log n), in place, not stable
//   - container/heap gives you the algorithms; you implement heap.Interface
package main

import (
	"cmp"
	"container/heap"
	"fmt"
	"math/rand/v2"
	"os"
	"slices"
)

// MinHeap is a from-scratch binary heap with a custom less function.
type MinHeap[T any] struct {
	items []T
	less  func(a, b T) bool
}

func NewMinHeap[T any](less func(a, b T) bool) *MinHeap[T] { return &MinHeap[T]{less: less} }

func (h *MinHeap[T]) Len() int { return len(h.items) }

func (h *MinHeap[T]) Push(v T) {
	h.items = append(h.items, v)
	i := len(h.items) - 1
	for i > 0 { // sift up
		parent := (i - 1) / 2
		if !h.less(h.items[i], h.items[parent]) {
			break
		}
		h.items[i], h.items[parent] = h.items[parent], h.items[i]
		i = parent
	}
}

func (h *MinHeap[T]) Pop() (T, bool) {
	var zero T
	if len(h.items) == 0 {
		return zero, false
	}
	top := h.items[0]
	last := len(h.items) - 1
	h.items[0] = h.items[last]
	h.items = h.items[:last]
	h.siftDown(0)
	return top, true
}

func (h *MinHeap[T]) siftDown(i int) {
	n := len(h.items)
	for {
		smallest := i
		l, r := 2*i+1, 2*i+2
		if l < n && h.less(h.items[l], h.items[smallest]) {
			smallest = l
		}
		if r < n && h.less(h.items[r], h.items[smallest]) {
			smallest = r
		}
		if smallest == i {
			return
		}
		h.items[i], h.items[smallest] = h.items[smallest], h.items[i]
		i = smallest
	}
}

// TopK returns the k largest values, largest first, using a min-heap of size k.
func TopK[T cmp.Ordered](values []T, k int) []T {
	h := NewMinHeap(func(a, b T) bool { return a < b })
	for _, v := range values {
		h.Push(v)
		if h.Len() > k {
			h.Pop() // drop the smallest of the current top-k
		}
	}
	out := make([]T, 0, k)
	for h.Len() > 0 {
		v, _ := h.Pop()
		out = append(out, v)
	}
	slices.Reverse(out)
	return out
}

type task struct {
	name     string
	priority int
}

// taskQueue implements heap.Interface for the standard library version.
type taskQueue []task

func (q taskQueue) Len() int           { return len(q) }
func (q taskQueue) Less(i, j int) bool { return q[i].priority < q[j].priority }
func (q taskQueue) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }
func (q *taskQueue) Push(x any)        { *q = append(*q, x.(task)) }
func (q *taskQueue) Pop() any {
	old := *q
	t := old[len(old)-1]
	*q = old[:len(old)-1]
	return t
}

func main() {
	fmt.Println("Topic: min-heap and priority queue")
	failed := false

	// Heap sort with our heap, checked against slices.Sort
	r := rand.New(rand.NewPCG(9, 9))
	data := make([]int, 2000)
	for i := range data {
		data[i] = r.IntN(10_000)
	}
	h := NewMinHeap(func(a, b int) bool { return a < b })
	for _, v := range data {
		h.Push(v)
	}
	var popped []int
	for h.Len() > 0 {
		v, _ := h.Pop()
		popped = append(popped, v)
	}
	want := slices.Clone(data)
	slices.Sort(want)
	fmt.Println("heap sort of 2000 items matches slices.Sort:", slices.Equal(popped, want))
	if !slices.Equal(popped, want) {
		failed = true
	}

	// Priority queue with our heap
	pq := NewMinHeap(func(a, b task) bool { return a.priority < b.priority })
	pq.Push(task{"write report", 3})
	pq.Push(task{"fix prod bug", 1})
	pq.Push(task{"lunch", 5})
	pq.Push(task{"code review", 2})
	fmt.Print("tasks by priority: ")
	for pq.Len() > 0 {
		t, _ := pq.Pop()
		fmt.Printf("%s(%d) ", t.name, t.priority)
	}
	fmt.Println()

	// Top-k
	scores := []int{55, 91, 12, 78, 91, 34, 66, 99, 40}
	top3 := TopK(scores, 3)
	fmt.Println("top 3 of", scores, "=", top3)
	if !slices.Equal(top3, []int{99, 91, 91}) {
		failed = true
	}

	// Standard library container/heap
	q := &taskQueue{{"b", 2}, {"c", 3}, {"a", 1}}
	heap.Init(q) // O(n) heapify
	heap.Push(q, task{"urgent", 0})
	first := heap.Pop(q).(task)
	fmt.Println("container/heap first pop:", first.name)
	if first.name != "urgent" {
		failed = true
	}

	if failed {
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}
