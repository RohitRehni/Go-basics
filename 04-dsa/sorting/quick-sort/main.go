// Topic: quick-sort | Concept: sorting | Stage: 04-dsa
// Run:   go run ./04-dsa/sorting/quick-sort
//
// Explore:
// * Pick a pivot, partition: smaller items left, larger right, recurse on both sides
// * Time O(n log n) average, O(n^2) worst (bad pivots on sorted input)
// * Space O(log n) for recursion on average; sorts in place
// * Lomuto partition (simple) vs Hoare partition (fewer swaps)
// * Random or median-of-three pivots avoid the sorted-input worst case
// * Not stable. Go's slices.Sort is pdqsort: quicksort + heapsort + insertion sort
package main

import (
	"cmp"
	"fmt"
	"math/rand/v2"
	"os"
	"slices"
)

// QuickSort sorts a in place using Lomuto partitioning and a random pivot.
func QuickSort[T cmp.Ordered](a []T) {
	if len(a) <= 1 {
		return
	}
	p := partition(a)
	QuickSort(a[:p])
	QuickSort(a[p+1:])
}

// partition moves a random pivot to its final position and returns that index.
func partition[T cmp.Ordered](a []T) int {
	r := rand.IntN(len(a))
	a[r], a[len(a)-1] = a[len(a)-1], a[r] // pivot to the end
	pivot := a[len(a)-1]
	store := 0
	for i := 0; i < len(a)-1; i++ {
		if a[i] < pivot {
			a[i], a[store] = a[store], a[i]
			store++
		}
	}
	a[store], a[len(a)-1] = a[len(a)-1], a[store]
	return store
}

// QuickSelect returns the k-th smallest item (0-based) in O(n) average time.
func QuickSelect[T cmp.Ordered](a []T, k int) T {
	a = slices.Clone(a)
	for {
		p := partition(a)
		switch {
		case p == k:
			return a[p]
		case k < p:
			a = a[:p]
		default:
			k -= p + 1
			a = a[p+1:]
		}
	}
}

func main() {
	fmt.Println("Topic: quick sort")

	failed := false
	demo := []int{10, 80, 30, 90, 40, 50, 70}
	sorted := slices.Clone(demo)
	QuickSort(sorted)
	fmt.Println("input: ", demo)
	fmt.Println("sorted:", sorted)

	r := rand.New(rand.NewPCG(7, 7))
	cases := map[string][]int{
		"random":       nil,
		"sorted":       nil,
		"reverse":      nil,
		"all equal":    nil,
		"many repeats": nil,
	}
	const n = 20_000
	for name := range cases {
		data := make([]int, n)
		for i := range data {
			switch name {
			case "random":
				data[i] = r.IntN(1_000_000)
			case "sorted":
				data[i] = i
			case "reverse":
				data[i] = n - i
			case "all equal":
				data[i] = 5
			case "many repeats":
				data[i] = r.IntN(10)
			}
		}
		cases[name] = data
	}
	for _, name := range []string{"random", "sorted", "reverse", "all equal", "many repeats"} {
		got := slices.Clone(cases[name])
		QuickSort(got)
		status := "ok"
		if !slices.IsSorted(got) {
			status = "MISMATCH"
			failed = true
		}
		fmt.Printf("  %-13s n=%d sorted correctly: %s\n", name, n, status)
	}

	nums := []int{7, 2, 9, 4, 1, 8}
	want := slices.Clone(nums)
	slices.Sort(want)
	for k := range nums {
		if QuickSelect(nums, k) != want[k] {
			failed = true
		}
	}
	fmt.Printf("QuickSelect: 3rd smallest of %v = %d\n", nums, QuickSelect(nums, 2))

	if failed {
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}
