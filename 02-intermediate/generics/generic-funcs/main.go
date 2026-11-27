// Topic: generic-funcs | Concept: generics | Stage: 02-intermediate
// Run:   go run ./02-intermediate/generics/generic-funcs
//
// Explore:
// * Type parameters in square brackets: func Map[T, U any](...)
// * Type inference: Map(nums, f) usually needs no [int, string]
// * Generic types: type Stack[T any] struct{ items []T }
// * Methods on generic types use the receiver's type parameter
// * `comparable` lets you use == and map keys
package main

import (
	"fmt"
	"strings"
)

func Map[T, U any](items []T, f func(T) U) []U {
	out := make([]U, 0, len(items))
	for _, it := range items {
		out = append(out, f(it))
	}
	return out
}

func Filter[T any](items []T, keep func(T) bool) []T {
	var out []T
	for _, it := range items {
		if keep(it) {
			out = append(out, it)
		}
	}
	return out
}

func Contains[T comparable](items []T, target T) bool {
	for _, it := range items {
		if it == target {
			return true
		}
	}
	return false
}

type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(v T) {
	s.items = append(s.items, v)
}

func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	v := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return v, true
}

func main() {
	fmt.Println("Topic: generic functions and types")

	nums := []int{1, 2, 3, 4, 5}
	squares := Map(nums, func(n int) int { return n * n })
	labels := Map(nums, func(n int) string { return fmt.Sprintf("#%d", n) })
	fmt.Println("squares:", squares)
	fmt.Println("labels:", labels)

	evens := Filter(nums, func(n int) bool { return n%2 == 0 })
	fmt.Println("evens:", evens)

	words := []string{"go", "rust", "python"}
	fmt.Println("upper:", Map(words, strings.ToUpper))
	fmt.Println("Contains(words, \"rust\"):", Contains(words, "rust"))

	var s Stack[string]
	s.Push("a")
	s.Push("b")
	top, _ := s.Pop()
	fmt.Println("popped:", top)
	_, _ = s.Pop()
	_, ok := s.Pop()
	fmt.Println("pop from empty stack ok:", ok)

	// Contains([]func(){}, nil) // compile error: func() does not satisfy comparable
}
