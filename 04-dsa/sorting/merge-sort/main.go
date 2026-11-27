// Topic: merge-sort | Concept: sorting | Stage: 04-dsa
// Run:   go run ./04-dsa/sorting/merge-sort
//
// Explore:
// * Divide and conquer: split in half, sort each half, merge the two sorted halves
// * Time O(n log n) in every case; space O(n) for the merge buffer
// * Stable: on ties, take from the left half first
// * Merging two sorted lists is a useful skill on its own (merge k logs, etc.)
// * Bottom-up version: merge runs of 1, 2, 4, ... with no recursion
package main

import (
	"cmp"
	"fmt"
	"math/rand/v2"
	"os"
	"slices"
)

// MergeSort returns a new sorted slice (top-down, recursive).
func MergeSort[T cmp.Ordered](a []T) []T {
	if len(a) <= 1 {
		return slices.Clone(a)
	}
	mid := len(a) / 2
	return merge(MergeSort(a[:mid]), MergeSort(a[mid:]))
}

// merge combines two sorted slices into one sorted slice.
func merge[T cmp.Ordered](left, right []T) []T {
	out := make([]T, 0, len(left)+len(right))
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		if left[i] <= right[j] { // <= keeps it stable
			out = append(out, left[i])
			i++
		} else {
			out = append(out, right[j])
			j++
		}
	}
	out = append(out, left[i:]...)
	return append(out, right[j:]...)
}

// MergeSortBottomUp sorts in passes of width 1, 2, 4, ... without recursion.
func MergeSortBottomUp[T cmp.Ordered](a []T) []T {
	src := slices.Clone(a)
	for width := 1; width < len(src); width *= 2 {
		dst := make([]T, 0, len(src))
		for lo := 0; lo < len(src); lo += 2 * width {
			mid := min(lo+width, len(src))
			hi := min(lo+2*width, len(src))
			dst = append(dst, merge(src[lo:mid], src[mid:hi])...)
		}
		src = dst
	}
	return src
}

func main() {
	fmt.Println("Topic: merge sort")

	failed := false
	small := []int{38, 27, 43, 3, 9, 82, 10}
	fmt.Println("input:     ", small)
	fmt.Println("top-down:  ", MergeSort(small))
	fmt.Println("bottom-up: ", MergeSortBottomUp(small))
	fmt.Println("merge([1 4 9], [2 3 10]) =", merge([]int{1, 4, 9}, []int{2, 3, 10}))

	r := rand.New(rand.NewPCG(1, 2)) // fixed seed: same numbers every run
	for _, n := range []int{0, 1, 2, 17, 1000, 50_000} {
		data := make([]int, n)
		for i := range data {
			data[i] = r.IntN(1000)
		}
		want := slices.Clone(data)
		slices.Sort(want)
		ok := slices.Equal(MergeSort(data), want) && slices.Equal(MergeSortBottomUp(data), want)
		status := "ok"
		if !ok {
			status = "MISMATCH"
			failed = true
		}
		fmt.Printf("  random n=%-6d both versions match slices.Sort: %s\n", n, status)
	}

	if failed {
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}
