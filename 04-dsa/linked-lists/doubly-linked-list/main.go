// Topic: doubly-linked-list | Concept: linked-lists | Stage: 04-dsa
// Run:   go run ./04-dsa/linked-lists/doubly-linked-list
//
// Explore:
//   - Each node points to next AND prev: walk both ways, remove a known node in O(1)
//   - Sentinel (dummy) node: head and tail are sentinel.next / sentinel.prev, so
//     insert/remove never needs nil checks
//   - Costs one extra pointer per node compared with a singly linked list
//   - Used inside LRU caches: a map finds the node, the list keeps recency order
//   - container/list in the standard library is a doubly linked list
package main

import (
	"container/list"
	"fmt"
	"os"
	"slices"
)

type dnode[T any] struct {
	val        T
	prev, next *dnode[T]
}

// DList is a circular doubly linked list with a sentinel node.
type DList[T any] struct {
	root dnode[T] // root.next = first, root.prev = last
	size int
}

func NewDList[T any]() *DList[T] {
	l := &DList[T]{}
	l.root.next, l.root.prev = &l.root, &l.root
	return l
}

func (l *DList[T]) insertAfter(at *dnode[T], v T) *dnode[T] {
	n := &dnode[T]{val: v, prev: at, next: at.next}
	at.next.prev = n
	at.next = n
	l.size++
	return n
}

func (l *DList[T]) PushFront(v T) *dnode[T] { return l.insertAfter(&l.root, v) }
func (l *DList[T]) PushBack(v T) *dnode[T]  { return l.insertAfter(l.root.prev, v) }

// Remove unlinks n in O(1): no search needed because n knows its neighbours.
func (l *DList[T]) Remove(n *dnode[T]) T {
	n.prev.next = n.next
	n.next.prev = n.prev
	n.prev, n.next = nil, nil
	l.size--
	return n.val
}

// MoveToFront relinks n right after the sentinel in O(1), without allocating.
func (l *DList[T]) MoveToFront(n *dnode[T]) {
	if l.root.next == n {
		return
	}
	n.prev.next = n.next // unlink
	n.next.prev = n.prev
	n.prev = &l.root // relink at the front
	n.next = l.root.next
	l.root.next.prev = n
	l.root.next = n
}

func (l *DList[T]) Forward() []T {
	var out []T
	for n := l.root.next; n != &l.root; n = n.next {
		out = append(out, n.val)
	}
	return out
}

func (l *DList[T]) Backward() []T {
	var out []T
	for n := l.root.prev; n != &l.root; n = n.prev {
		out = append(out, n.val)
	}
	return out
}

var failed bool

func expect(label string, got, want []string) {
	status := "ok"
	if !slices.Equal(got, want) {
		status = fmt.Sprintf("MISMATCH (want %v)", want)
		failed = true
	}
	fmt.Printf("  %-24s %v %s\n", label, got, status)
}

func main() {
	fmt.Println("Topic: doubly linked list")

	l := NewDList[string]()
	b := l.PushBack("b")
	l.PushBack("c")
	l.PushFront("a")
	d := l.PushBack("d")
	expect("forward", l.Forward(), []string{"a", "b", "c", "d"})
	expect("backward", l.Backward(), []string{"d", "c", "b", "a"})

	l.Remove(b)
	expect("remove b (O(1))", l.Forward(), []string{"a", "c", "d"})

	l.MoveToFront(d)
	expect("move d to front", l.Forward(), []string{"d", "a", "c"})
	expect("backward again", l.Backward(), []string{"c", "a", "d"})
	fmt.Println("  size:", l.size)
	if l.size != 3 {
		failed = true
	}

	// The standard library version
	std := list.New()
	std.PushBack(2)
	std.PushFront(1)
	std.PushBack(3)
	fmt.Print("container/list forward: ")
	for e := std.Front(); e != nil; e = e.Next() {
		fmt.Print(e.Value, " ")
	}
	fmt.Println()

	if failed {
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}
