// Topic: dijkstra | Concept: graphs | Stage: 04-dsa
// Run:   go run ./04-dsa/graphs/dijkstra
//
// Explore:
// * Shortest paths from one source in a graph with NON-NEGATIVE edge weights
// * Greedy: always settle the unvisited node with the smallest known distance
// * With a min-heap: O((V + E) log V); with a plain array scan: O(V^2)
// * "Lazy deletion": push duplicates into the heap and skip stale entries when popped
// * Keep prev[node] to rebuild the path
// * Negative weights break it: use Bellman-Ford instead
package main

import (
	"container/heap"
	"fmt"
	"math"
	"os"
	"slices"
)

type Edge struct {
	to     string
	weight int
}

type Graph map[string][]Edge

func (g Graph) AddRoad(a, b string, w int) { // undirected
	g[a] = append(g[a], Edge{b, w})
	g[b] = append(g[b], Edge{a, w})
}

type item struct {
	node string
	dist int
}
type pq []item

func (q pq) Len() int           { return len(q) }
func (q pq) Less(i, j int) bool { return q[i].dist < q[j].dist }
func (q pq) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }
func (q *pq) Push(x any)        { *q = append(*q, x.(item)) }
func (q *pq) Pop() any {
	old := *q
	it := old[len(old)-1]
	*q = old[:len(old)-1]
	return it
}

// Dijkstra returns the distance to every reachable node and the previous hop on its path.
func Dijkstra(g Graph, src string) (map[string]int, map[string]string) {
	dist := map[string]int{src: 0}
	prev := map[string]string{}
	q := &pq{{src, 0}}
	for q.Len() > 0 {
		cur := heap.Pop(q).(item)
		if cur.dist > dist[cur.node] {
			continue // stale entry: a shorter path was already found
		}
		for _, e := range g[cur.node] {
			nd := cur.dist + e.weight
			if old, seen := dist[e.to]; !seen || nd < old {
				dist[e.to] = nd
				prev[e.to] = cur.node
				heap.Push(q, item{e.to, nd})
			}
		}
	}
	return dist, prev
}

func pathTo(prev map[string]string, src, dst string) []string {
	path := []string{dst}
	for n := dst; n != src; {
		p, ok := prev[n]
		if !ok {
			return nil
		}
		path = append(path, p)
		n = p
	}
	slices.Reverse(path)
	return path
}

// bellmanFord is a slow O(V*E) reference used to double-check Dijkstra.
func bellmanFord(g Graph, src string) map[string]int {
	dist := map[string]int{}
	for n := range g {
		dist[n] = math.MaxInt
	}
	dist[src] = 0
	for range len(g) - 1 {
		for a, edges := range g {
			if dist[a] == math.MaxInt {
				continue
			}
			for _, e := range edges {
				dist[e.to] = min(dist[e.to], dist[a]+e.weight)
			}
		}
	}
	return dist
}

func main() {
	fmt.Println("Topic: Dijkstra's shortest path")
	failed := false

	g := Graph{}
	g.AddRoad("Home", "Bakery", 4)
	g.AddRoad("Home", "Park", 1)
	g.AddRoad("Park", "Bakery", 2)
	g.AddRoad("Park", "School", 5)
	g.AddRoad("Bakery", "School", 1)
	g.AddRoad("School", "Office", 3)
	g.AddRoad("Bakery", "Office", 7)

	dist, prev := Dijkstra(g, "Home")
	ref := bellmanFord(g, "Home")
	for _, dst := range []string{"Park", "Bakery", "School", "Office"} {
		status := "ok"
		if dist[dst] != ref[dst] {
			status = fmt.Sprintf("MISMATCH (Bellman-Ford says %d)", ref[dst])
			failed = true
		}
		fmt.Printf("  Home -> %-7s distance %2d via %v %s\n", dst, dist[dst], pathTo(prev, "Home", dst), status)
	}
	if dist["Office"] != 7 {
		failed = true
	}

	if failed {
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}
