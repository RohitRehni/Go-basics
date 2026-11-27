// Topic: topological-sort | Concept: graphs | Stage: 04-dsa
// Run:   go run ./04-dsa/graphs/topological-sort
//
// Explore:
//   - Orders the nodes of a DIRECTED graph so every edge a -> b has a before b
//   - Only possible when there is no cycle (a DAG: directed acyclic graph)
//   - Kahn's algorithm: repeatedly take nodes with in-degree 0 (a queue); if some
//     nodes are never taken, there is a cycle
//   - DFS version: add a node after all its descendants, then reverse
//   - Both are O(V + E). Uses: build order, course prerequisites, task scheduling, DB migrations
package main

import (
	"fmt"
	"os"
	"slices"
)

type DAG map[string][]string // node -> nodes that depend on it

// Kahn returns a topological order, or an error naming the nodes stuck in a cycle.
func Kahn(g DAG) ([]string, error) {
	indeg := map[string]int{}
	for n, outs := range g {
		if _, ok := indeg[n]; !ok {
			indeg[n] = 0
		}
		for _, m := range outs {
			indeg[m]++
		}
	}
	var ready []string
	for n, d := range indeg {
		if d == 0 {
			ready = append(ready, n)
		}
	}
	slices.Sort(ready) // deterministic output
	var order []string
	for len(ready) > 0 {
		n := ready[0]
		ready = ready[1:]
		order = append(order, n)
		next := slices.Clone(g[n])
		slices.Sort(next)
		for _, m := range next {
			indeg[m]--
			if indeg[m] == 0 {
				ready = append(ready, m)
			}
		}
	}
	if len(order) != len(indeg) {
		var stuck []string
		for n, d := range indeg {
			if d > 0 {
				stuck = append(stuck, n)
			}
		}
		slices.Sort(stuck)
		return order, fmt.Errorf("cycle detected among %v", stuck)
	}
	return order, nil
}

// valid checks that every edge goes forward in order.
func valid(g DAG, order []string) bool {
	pos := map[string]int{}
	for i, n := range order {
		pos[n] = i
	}
	for a, outs := range g {
		for _, b := range outs {
			if pos[a] >= pos[b] {
				return false
			}
		}
	}
	return true
}

func main() {
	fmt.Println("Topic: topological sort (Kahn's algorithm)")
	failed := false

	// Course prerequisites: "intro" must come before "data-structures", etc.
	courses := DAG{
		"intro":           {"data-structures", "web-basics"},
		"data-structures": {"algorithms", "databases"},
		"web-basics":      {"backend"},
		"databases":       {"backend"},
		"algorithms":      {},
		"backend":         {"capstone"},
		"capstone":        {},
	}
	order, err := Kahn(courses)
	fmt.Println("course order:", order, "| error:", err)
	fmt.Println("every prerequisite comes first:", valid(courses, order))
	if err != nil || !valid(courses, order) {
		failed = true
	}

	cyclic := DAG{"a": {"b"}, "b": {"c"}, "c": {"a"}, "start": {"a"}}
	_, err = Kahn(cyclic)
	fmt.Println("cyclic graph:", err)
	if err == nil {
		failed = true
	}

	if failed {
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}
