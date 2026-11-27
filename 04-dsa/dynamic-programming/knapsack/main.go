// Topic: knapsack | Concept: dynamic-programming | Stage: 04-dsa
// Run:   go run ./04-dsa/dynamic-programming/knapsack
//
// Explore:
//   - 0/1 knapsack: pick items (each at most once) to maximise value within a weight limit
//   - dp[i][w] = best value using the first i items with capacity w:
//     max(skip item i, take item i if it fits)
//   - Time O(n * W), space O(n * W); a 1-D table iterated backwards gives O(W) space
//   - Walk the table backwards to find WHICH items were chosen
//   - Greedy by value/weight ratio is NOT optimal for 0/1 knapsack (shown below)
//   - Brute force tries all 2^n subsets: fine for checking small cases
package main

import (
	"fmt"
	"os"
	"slices"
)

type Item struct {
	Name   string
	Weight int
	Value  int
}

// Knapsack returns the best value and the chosen item names.
func Knapsack(items []Item, capacity int) (int, []string) {
	n := len(items)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, capacity+1)
	}
	for i := 1; i <= n; i++ {
		it := items[i-1]
		for w := 0; w <= capacity; w++ {
			dp[i][w] = dp[i-1][w] // skip
			if it.Weight <= w {
				dp[i][w] = max(dp[i][w], dp[i-1][w-it.Weight]+it.Value) // take
			}
		}
	}
	var chosen []string
	w := capacity
	for i := n; i > 0; i-- {
		if dp[i][w] != dp[i-1][w] { // value changed, so item i was taken
			chosen = append(chosen, items[i-1].Name)
			w -= items[i-1].Weight
		}
	}
	slices.Reverse(chosen)
	return dp[n][capacity], chosen
}

// KnapsackSmall uses a single row, iterating capacity downwards.
func KnapsackSmall(items []Item, capacity int) int {
	dp := make([]int, capacity+1)
	for _, it := range items {
		for w := capacity; w >= it.Weight; w-- {
			dp[w] = max(dp[w], dp[w-it.Weight]+it.Value)
		}
	}
	return dp[capacity]
}

// bruteForce checks every subset (2^n).
func bruteForce(items []Item, capacity int) int {
	best := 0
	for mask := 0; mask < 1<<len(items); mask++ {
		w, v := 0, 0
		for i, it := range items {
			if mask&(1<<i) != 0 {
				w += it.Weight
				v += it.Value
			}
		}
		if w <= capacity {
			best = max(best, v)
		}
	}
	return best
}

func greedyByRatio(items []Item, capacity int) int {
	sorted := slices.Clone(items)
	slices.SortFunc(sorted, func(a, b Item) int {
		return b.Value*a.Weight - a.Value*b.Weight // higher value/weight first
	})
	v := 0
	for _, it := range sorted {
		if it.Weight <= capacity {
			capacity -= it.Weight
			v += it.Value
		}
	}
	return v
}

func main() {
	fmt.Println("Topic: 0/1 knapsack")
	failed := false

	// Classic counter-example: greedy grabs the best ratio (camera) and gets stuck.
	items := []Item{
		{"camera", 1, 60},
		{"laptop", 2, 100},
		{"guitar", 3, 120},
	}
	capacity := 5

	best, chosen := Knapsack(items, capacity)
	fmt.Printf("capacity %d: best value %d with %v\n", capacity, best, chosen)
	fmt.Println("one-row version:", KnapsackSmall(items, capacity))
	fmt.Println("brute force (all 2^3 subsets):", bruteForce(items, capacity))
	greedy := greedyByRatio(items, capacity)
	fmt.Println("greedy by value/weight:", greedy, "(worse than the DP answer)")
	if best != bruteForce(items, capacity) || best != KnapsackSmall(items, capacity) || best != 220 || greedy >= best {
		failed = true
	}

	// A bigger set, checked against brute force for every capacity
	items = append(items, Item{"books", 2, 30}, Item{"stereo", 4, 130}, Item{"tent", 5, 160})

	allMatch := true
	for c := range 18 {
		if b, _ := Knapsack(items, c); b != bruteForce(items, c) {
			fmt.Println("MISMATCH at capacity", c)
			allMatch = false
		}
	}
	fmt.Println("6 items: DP matches brute force for every capacity 0..17:", allMatch)
	if !allMatch {
		failed = true
	}

	if failed {
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}
