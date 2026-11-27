// Topic: binary-search | Concept: searching | Stage: 04-dsa
// Run:   go run ./04-dsa/searching/binary-search
//
// Explore:
// * Needs SORTED input; halve the search range each step
// * Time O(log n): 1,000,000 items need at most 20 comparisons; space O(1)
// * mid := lo + (hi-lo)/2 avoids integer overflow on huge slices
// * Lower bound: first index where items[i] >= target (insert position)
// * Recursive version uses O(log n) stack; the loop version uses O(1)
// * slices.BinarySearch / sort.SearchInts are the standard-library versions
package main

import (
	"cmp"
	"fmt"
	"os"
	"slices"
)

// BinarySearch returns the index of target in sorted items, or -1.
func BinarySearch[T cmp.Ordered](items []T, target T) int {
	lo, hi := 0, len(items)-1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		switch {
		case items[mid] == target:
			return mid
		case items[mid] < target:
			lo = mid + 1
		default:
			hi = mid - 1
		}
	}
	return -1
}

// LowerBound returns the first index i with items[i] >= target.
func LowerBound[T cmp.Ordered](items []T, target T) int {
	lo, hi := 0, len(items)
	for lo < hi {
		mid := lo + (hi-lo)/2
		if items[mid] < target {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo
}

// BinarySearchRec is the recursive version.
func BinarySearchRec[T cmp.Ordered](items []T, target T, lo, hi int) int {
	if lo > hi {
		return -1
	}
	mid := lo + (hi-lo)/2
	if items[mid] == target {
		return mid
	}
	if items[mid] < target {
		return BinarySearchRec(items, target, mid+1, hi)
	}
	return BinarySearchRec(items, target, lo, mid-1)
}

var failed bool

func check(label string, got, want int) {
	status := "ok"
	if got != want {
		status = fmt.Sprintf("MISMATCH (want %d)", want)
		failed = true
	}
	fmt.Printf("  %-28s = %3d  %s\n", label, got, status)
}

func main() {
	fmt.Println("Topic: binary search")

	sorted := []int{2, 5, 8, 12, 16, 23, 38, 56, 72, 91}
	fmt.Println("sorted:", sorted)

	for _, target := range []int{2, 23, 91, 4, 100} {
		want, found := slices.BinarySearch(sorted, target)
		if !found {
			want = -1
		}
		check(fmt.Sprintf("BinarySearch(%d)", target), BinarySearch(sorted, target), want)
		check(fmt.Sprintf("BinarySearchRec(%d)", target), BinarySearchRec(sorted, target, 0, len(sorted)-1), want)
	}

	for _, target := range []int{1, 13, 91, 100} {
		want, _ := slices.BinarySearch(sorted, target) // insert position
		check(fmt.Sprintf("LowerBound(%d)", target), LowerBound(sorted, target), want)
	}

	// How many steps for a big slice?
	big := make([]int, 1_000_000)
	for i := range big {
		big[i] = i * 2
	}
	steps := 0
	lo, hi := 0, len(big)-1
	for lo <= hi {
		steps++
		mid := lo + (hi-lo)/2
		if big[mid] == 1_999_998 {
			break
		} else if big[mid] < 1_999_998 {
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	fmt.Printf("found the last of 1,000,000 items in %d steps\n", steps)

	if failed {
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}
