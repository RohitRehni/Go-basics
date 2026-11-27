// Topic: permutations-subsets | Concept: backtracking | Stage: 04-dsa
// Run:   go run ./04-dsa/backtracking/permutations-subsets
//
// Explore:
// * Backtracking = build a solution step by step, undo the last choice, try the next one
// * Template: choose -> recurse -> un-choose
// * Subsets: 2^n of them; each item is either in or out
// * Permutations: n! orderings; track which items are already used
// * Combinations (choose k of n) prune branches that cannot reach k items
// * Copy the current path when saving a result: the slice keeps changing
package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

func Subsets[T any](items []T) [][]T {
	var out [][]T
	var path []T
	var walk func(start int)
	walk = func(start int) {
		out = append(out, slices.Clone(path)) // every path is a subset
		for i := start; i < len(items); i++ {
			path = append(path, items[i]) // choose
			walk(i + 1)                   // explore
			path = path[:len(path)-1]     // un-choose
		}
	}
	walk(0)
	return out
}

func Permutations[T any](items []T) [][]T {
	var out [][]T
	used := make([]bool, len(items))
	var path []T
	var walk func()
	walk = func() {
		if len(path) == len(items) {
			out = append(out, slices.Clone(path))
			return
		}
		for i, v := range items {
			if used[i] {
				continue
			}
			used[i] = true
			path = append(path, v)
			walk()
			path = path[:len(path)-1]
			used[i] = false
		}
	}
	walk()
	return out
}

func Combinations(n, k int) [][]int {
	var out [][]int
	var path []int
	var walk func(start int)
	walk = func(start int) {
		if len(path) == k {
			out = append(out, slices.Clone(path))
			return
		}
		for i := start; i <= n-(k-len(path))+1; i++ { // prune: leave room for the rest
			path = append(path, i)
			walk(i + 1)
			path = path[:len(path)-1]
		}
	}
	walk(1)
	return out
}

// CombinationSum finds sets of candidates (reusable) that add up to target.
func CombinationSum(candidates []int, target int) [][]int {
	var out [][]int
	var path []int
	var walk func(start, remaining int)
	walk = func(start, remaining int) {
		if remaining == 0 {
			out = append(out, slices.Clone(path))
			return
		}
		for i := start; i < len(candidates); i++ {
			if candidates[i] > remaining {
				continue
			}
			path = append(path, candidates[i])
			walk(i, remaining-candidates[i]) // i, not i+1: the same number may repeat
			path = path[:len(path)-1]
		}
	}
	walk(0, target)
	return out
}

func main() {
	fmt.Println("Topic: backtracking - subsets, permutations, combinations")
	failed := false

	subs := Subsets([]string{"a", "b", "c"})
	fmt.Println("subsets of [a b c]:", subs)
	if len(subs) != 8 {
		failed = true
	}

	perms := Permutations([]int{1, 2, 3})
	fmt.Println("permutations of [1 2 3]:", perms)
	seen := map[string]bool{}
	for _, p := range perms {
		seen[fmt.Sprint(p)] = true
	}
	if len(perms) != 6 || len(seen) != 6 {
		failed = true
	}
	fmt.Println("count of permutations of 6 items:", len(Permutations([]int{1, 2, 3, 4, 5, 6})), "(6! = 720)")
	if len(Permutations([]int{1, 2, 3, 4, 5, 6})) != 720 {
		failed = true
	}

	combos := Combinations(5, 3)
	fmt.Println("choose 3 of 1..5:", len(combos), "combinations, first:", combos[0], "last:", combos[len(combos)-1])
	if len(combos) != 10 {
		failed = true
	}

	sums := CombinationSum([]int{2, 3, 6, 7}, 7)
	var lines []string
	for _, s := range sums {
		lines = append(lines, fmt.Sprint(s))
	}
	fmt.Println("combinations of [2 3 6 7] summing to 7:", strings.Join(lines, " "))
	if len(sums) != 2 {
		failed = true
	}

	if failed {
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}
