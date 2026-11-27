// Topic: simple-sorts | Concept: sorting | Stage: 04-dsa
// Run:   go run ./04-dsa/sorting/simple-sorts
//
// Explore:
//   - Bubble sort: swap neighbours until no swaps happen; O(n^2), O(n) if already sorted
//   - Selection sort: pick the minimum of the rest each round; always O(n^2), few swaps
//   - Insertion sort: grow a sorted prefix, slide each item into place;
//     O(n^2) worst, O(n) on nearly sorted data, which is why real sorts use it for small runs
//   - All three sort in place: space O(1)
//   - Stable = equal items keep their order: bubble and insertion yes, selection no
package main

import (
	"cmp"
	"fmt"
	"os"
	"slices"
)

func BubbleSort[T cmp.Ordered](a []T) {
	for end := len(a) - 1; end > 0; end-- {
		swapped := false
		for i := 0; i < end; i++ {
			if a[i] > a[i+1] {
				a[i], a[i+1] = a[i+1], a[i]
				swapped = true
			}
		}
		if !swapped { // already sorted: stop early
			return
		}
	}
}

func SelectionSort[T cmp.Ordered](a []T) {
	for i := range a {
		minIdx := i
		for j := i + 1; j < len(a); j++ {
			if a[j] < a[minIdx] {
				minIdx = j
			}
		}
		a[i], a[minIdx] = a[minIdx], a[i]
	}
}

func InsertionSort[T cmp.Ordered](a []T) {
	for i := 1; i < len(a); i++ {
		key := a[i]
		j := i - 1
		for j >= 0 && a[j] > key {
			a[j+1] = a[j] // shift right
			j--
		}
		a[j+1] = key
	}
}

func main() {
	fmt.Println("Topic: simple sorts (bubble, selection, insertion)")

	inputs := [][]int{
		{29, 10, 14, 37, 13},
		{5, 4, 3, 2, 1},
		{1, 2, 3, 4, 5},
		{7, 7, 3, 7, 1},
		{},
		{42},
	}
	sorts := []struct {
		name string
		fn   func([]int)
	}{
		{"bubble", BubbleSort[int]},
		{"selection", SelectionSort[int]},
		{"insertion", InsertionSort[int]},
	}

	failed := false
	for _, in := range inputs {
		want := slices.Clone(in)
		slices.Sort(want)
		for _, s := range sorts {
			got := slices.Clone(in)
			s.fn(got)
			status := "ok"
			if !slices.Equal(got, want) {
				status = "MISMATCH"
				failed = true
			}
			fmt.Printf("  %-9s %-18v -> %-18v %s\n", s.name, in, got, status)
		}
	}

	words := []string{"pear", "apple", "fig", "banana"}
	InsertionSort(words)
	fmt.Println("insertion sort on strings:", words)
	if !slices.IsSorted(words) {
		failed = true
	}

	if failed {
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}
