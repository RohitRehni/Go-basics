// Topic: singly-linked-list | Concept: linked-lists | Stage: 04-dsa
// Run:   go run ./04-dsa/linked-lists/singly-linked-list
//
// Explore:
// * Each node holds a value and a pointer to the next node; the list keeps head (and tail)
// * Push front O(1); push back O(1) with a tail pointer, O(n) without
// * Find / delete by value O(n): you have to walk from the head
// * No random access: list[i] is O(n), unlike a slice
// * Classic exercises: reverse in place, find the middle (slow/fast pointers), detect a cycle
// * In real Go code a slice is almost always faster (cache-friendly); container/list exists too
package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

type node[T comparable] struct {
	val  T
	next *node[T]
}

type LinkedList[T comparable] struct {
	head, tail *node[T]
	size       int
}

func (l *LinkedList[T]) PushFront(v T) {
	l.head = &node[T]{val: v, next: l.head}
	if l.tail == nil {
		l.tail = l.head
	}
	l.size++
}

func (l *LinkedList[T]) PushBack(v T) {
	n := &node[T]{val: v}
	if l.tail == nil {
		l.head, l.tail = n, n
	} else {
		l.tail.next = n
		l.tail = n
	}
	l.size++
}

// Delete removes the first node with value v and reports whether it was found.
func (l *LinkedList[T]) Delete(v T) bool {
	var prev *node[T]
	for cur := l.head; cur != nil; prev, cur = cur, cur.next {
		if cur.val != v {
			continue
		}
		if prev == nil {
			l.head = cur.next
		} else {
			prev.next = cur.next
		}
		if cur == l.tail {
			l.tail = prev
		}
		l.size--
		return true
	}
	return false
}

func (l *LinkedList[T]) Contains(v T) bool {
	for cur := l.head; cur != nil; cur = cur.next {
		if cur.val == v {
			return true
		}
	}
	return false
}

// Reverse flips the list in place in O(n) time and O(1) space.
func (l *LinkedList[T]) Reverse() {
	var prev *node[T]
	cur := l.head
	l.tail = l.head
	for cur != nil {
		next := cur.next
		cur.next = prev
		prev, cur = cur, next
	}
	l.head = prev
}

// Middle uses slow/fast pointers: fast moves 2 steps for every 1 of slow.
func (l *LinkedList[T]) Middle() (T, bool) {
	var zero T
	if l.head == nil {
		return zero, false
	}
	slow, fast := l.head, l.head
	for fast.next != nil && fast.next.next != nil {
		slow = slow.next
		fast = fast.next.next
	}
	return slow.val, true
}

func (l *LinkedList[T]) Values() []T {
	out := make([]T, 0, l.size)
	for cur := l.head; cur != nil; cur = cur.next {
		out = append(out, cur.val)
	}
	return out
}

func (l *LinkedList[T]) String() string {
	var b strings.Builder
	for cur := l.head; cur != nil; cur = cur.next {
		fmt.Fprintf(&b, "%v -> ", cur.val)
	}
	b.WriteString("nil")
	return b.String()
}

var failed bool

func expect(label string, got, want []int) {
	status := "ok"
	if !slices.Equal(got, want) {
		status = fmt.Sprintf("MISMATCH (want %v)", want)
		failed = true
	}
	fmt.Printf("  %-22s %v %s\n", label, got, status)
}

func main() {
	fmt.Println("Topic: singly linked list")

	var l LinkedList[int]
	for _, v := range []int{2, 3, 4} {
		l.PushBack(v)
	}
	l.PushFront(1)
	fmt.Println("list:", &l)
	expect("after pushes", l.Values(), []int{1, 2, 3, 4})

	l.Delete(3)
	expect("delete 3", l.Values(), []int{1, 2, 4})
	l.Delete(1)
	expect("delete head 1", l.Values(), []int{2, 4})
	l.Delete(4)
	expect("delete tail 4", l.Values(), []int{2})
	l.PushBack(9) // tail must still be right after deleting the old tail
	expect("push back 9", l.Values(), []int{2, 9})
	fmt.Println("  delete 100 (missing):", l.Delete(100))

	for _, v := range []int{10, 11, 12} {
		l.PushBack(v)
	}
	l.Reverse()
	expect("reverse", l.Values(), []int{12, 11, 10, 9, 2})
	l.PushBack(1)
	expect("push back after reverse", l.Values(), []int{12, 11, 10, 9, 2, 1})

	mid, _ := l.Middle()
	fmt.Println("  middle:", mid, "| contains 10:", l.Contains(10), "| size:", l.size)
	if mid != 10 || l.size != 6 {
		failed = true
	}

	if failed {
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}
