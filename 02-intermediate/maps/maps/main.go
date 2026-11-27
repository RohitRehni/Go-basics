// Topic: maps | Concept: maps | Stage: 02-intermediate
// Run:   go run ./02-intermediate/maps/maps
//
// Explore:
// * make(map[K]V) and map literals
// * Reading a missing key returns the zero value; use v, ok := m[k] to tell
// * delete(m, k) is safe even if k is missing
// * Iteration order is random on purpose: sort the keys for stable output
// * A nil map can be read but writing to it panics
// * Maps are reference-like: passing one to a function shares it
// * Map of slices, and struct values (copy, modify, store back)
// * The maps package: maps.Keys, maps.Clone
package main

import (
	"fmt"
	"maps"
	"slices"
)

func addBonus(scores map[string]int) {
	for k := range scores {
		scores[k] += 5
	}
}

type Stock struct {
	Qty int
}

func main() {
	fmt.Println("Topic: maps")

	scores := map[string]int{"asha": 90, "ravi": 72}
	scores["meera"] = 85

	fmt.Println("missing key gives zero value:", scores["nobody"])
	if v, ok := scores["ravi"]; ok {
		fmt.Println("ravi found:", v)
	}

	delete(scores, "nobody") // no-op, no panic
	addBonus(scores)

	keys := slices.Sorted(maps.Keys(scores))
	for _, k := range keys {
		fmt.Printf("%-6s %d\n", k, scores[k])
	}
	fmt.Println("len =", len(scores))

	// Counting words
	counts := make(map[string]int)
	for _, w := range []string{"go", "rust", "go", "python", "go"} {
		counts[w]++
	}
	fmt.Println("word counts (fmt prints maps sorted):", counts)

	// Map of slices
	groups := map[string][]string{}
	groups["backend"] = append(groups["backend"], "go", "perl")
	fmt.Println("groups:", groups)

	// Struct values cannot be changed in place
	inventory := map[string]Stock{"apple": {Qty: 3}}
	// inventory["apple"].Qty++ // compile error: cannot assign to struct field inventory["apple"].Qty in map
	item := inventory["apple"]
	item.Qty++
	inventory["apple"] = item
	fmt.Println("inventory:", inventory)

	clone := maps.Clone(scores)
	clone["asha"] = 0
	fmt.Println("original asha after editing clone:", scores["asha"])

	var nilMap map[string]int
	fmt.Println("read from nil map:", nilMap["x"])
	// nilMap["x"] = 1 // runtime panic: assignment to entry in nil map
}
