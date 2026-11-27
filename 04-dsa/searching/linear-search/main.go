// Topic: linear-search | Concept: searching | Stage: 04-dsa
// Run:   go run ./04-dsa/searching/linear-search
//
// Explore:
// * Check every element in order until a match: works on unsorted data
// * Time O(n) worst and average, O(1) best (first element); space O(1)
// * Generic version with `comparable` works for ints, strings, structs
// * Predicate version (find first match of a condition) is often more useful
// * slices.Index / slices.IndexFunc are the standard-library equivalents
package main

import (
	"fmt"
	"os"
	"slices"
)

// LinearSearch returns the index of target, or -1.
func LinearSearch[T comparable](items []T, target T) int {
	for i, v := range items {
		if v == target {
			return i
		}
	}
	return -1
}

// FindFirst returns the index of the first item where match is true, or -1.
func FindFirst[T any](items []T, match func(T) bool) int {
	for i, v := range items {
		if match(v) {
			return i
		}
	}
	return -1
}

var failed bool

func check(label string, got, want int) {
	status := "ok"
	if got != want {
		status = fmt.Sprintf("MISMATCH (want %d)", want)
		failed = true
	}
	fmt.Printf("  %-32s = %3d  %s\n", label, got, status)
}

func main() {
	fmt.Println("Topic: linear search")

	nums := []int{42, 7, 19, 3, 88, 7}
	fmt.Println("nums:", nums)
	for _, target := range []int{42, 7, 88, 100} {
		check(fmt.Sprintf("LinearSearch(nums, %d)", target), LinearSearch(nums, target), slices.Index(nums, target))
	}

	names := []string{"go", "rust", "python", "perl"}
	check(`LinearSearch(names, "python")`, LinearSearch(names, "python"), slices.Index(names, "python"))

	firstOver50 := FindFirst(nums, func(n int) bool { return n > 50 })
	check("FindFirst(n > 50)", firstOver50, slices.IndexFunc(nums, func(n int) bool { return n > 50 }))

	if failed {
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}
