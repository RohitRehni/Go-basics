// Topic: bfs-dfs | Concept: graphs | Stage: 04-dsa
// Run:   go run ./04-dsa/graphs/bfs-dfs
//
// Explore:
// * Adjacency list: map[node][]neighbours. Space O(V + E)
// * BFS uses a queue and explores level by level: shortest path in UNWEIGHTED graphs
// * DFS uses recursion (or a stack) and goes deep first: cycles, connected components, mazes
// * Both are O(V + E) time; always keep a `visited` set or you loop forever on cycles
// * Record each node's parent during BFS to rebuild the actual path
// * Grids are graphs too: neighbours are up/down/left/right
package main

import (
	"fmt"
	"os"
	"slices"
)

type Graph map[string][]string

func (g Graph) AddEdge(a, b string) { // undirected
	g[a] = append(g[a], b)
	g[b] = append(g[b], a)
}

// BFS returns the visit order from start.
func (g Graph) BFS(start string) []string {
	visited := map[string]bool{start: true}
	queue := []string{start}
	var order []string
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		order = append(order, cur)
		for _, nb := range g[cur] {
			if !visited[nb] {
				visited[nb] = true
				queue = append(queue, nb)
			}
		}
	}
	return order
}

// DFS returns the visit order from start (recursive).
func (g Graph) DFS(start string) []string {
	visited := map[string]bool{}
	var order []string
	var visit func(string)
	visit = func(n string) {
		visited[n] = true
		order = append(order, n)
		for _, nb := range g[n] {
			if !visited[nb] {
				visit(nb)
			}
		}
	}
	visit(start)
	return order
}

// ShortestPath returns the path with the fewest edges, or nil.
func (g Graph) ShortestPath(from, to string) []string {
	parent := map[string]string{from: ""}
	queue := []string{from}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur == to {
			var path []string
			for n := to; n != ""; n = parent[n] {
				path = append(path, n)
			}
			slices.Reverse(path)
			return path
		}
		for _, nb := range g[cur] {
			if _, seen := parent[nb]; !seen {
				parent[nb] = cur
				queue = append(queue, nb)
			}
		}
	}
	return nil
}

// Components counts connected components with DFS.
func (g Graph) Components() int {
	visited := map[string]bool{}
	nodes := make([]string, 0, len(g))
	for n := range g {
		nodes = append(nodes, n)
	}
	slices.Sort(nodes)
	count := 0
	for _, n := range nodes {
		if visited[n] {
			continue
		}
		count++
		stack := []string{n} // iterative DFS
		for len(stack) > 0 {
			cur := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if visited[cur] {
				continue
			}
			visited[cur] = true
			stack = append(stack, g[cur]...)
		}
	}
	return count
}

// gridShortest finds the fewest steps through a maze ('#' = wall) with BFS.
func gridShortest(grid []string) int {
	type pt struct{ r, c int }
	rows, cols := len(grid), len(grid[0])
	dist := map[pt]int{{0, 0}: 0}
	queue := []pt{{0, 0}}
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		if p.r == rows-1 && p.c == cols-1 {
			return dist[p]
		}
		for _, d := range []pt{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			n := pt{p.r + d.r, p.c + d.c}
			if n.r < 0 || n.c < 0 || n.r >= rows || n.c >= cols || grid[n.r][n.c] == '#' {
				continue
			}
			if _, seen := dist[n]; !seen {
				dist[n] = dist[p] + 1
				queue = append(queue, n)
			}
		}
	}
	return -1
}

func main() {
	fmt.Println("Topic: BFS and DFS")
	failed := false

	g := Graph{}
	for _, e := range [][2]string{{"A", "B"}, {"A", "C"}, {"B", "D"}, {"C", "D"}, {"D", "E"}, {"E", "F"}, {"X", "Y"}} {
		g.AddEdge(e[0], e[1])
	}

	bfs := g.BFS("A")
	dfs := g.DFS("A")
	fmt.Println("BFS from A:", bfs)
	fmt.Println("DFS from A:", dfs)
	if !slices.Equal(bfs, []string{"A", "B", "C", "D", "E", "F"}) || !slices.Equal(dfs, []string{"A", "B", "D", "C", "E", "F"}) {
		failed = true
	}

	path := g.ShortestPath("A", "F")
	fmt.Println("shortest path A -> F:", path, "(", len(path)-1, "edges )")
	if len(path) != 5 {
		failed = true
	}
	fmt.Println("path A -> X:", g.ShortestPath("A", "X"), "(different component)")

	comps := g.Components()
	fmt.Println("connected components:", comps)
	if comps != 2 {
		failed = true
	}

	maze := []string{
		"..#....",
		".##.##.",
		"....#..",
		"#.#...#",
		"..#.#..",
	}
	steps := gridShortest(maze)
	fmt.Println("maze shortest path steps:", steps)
	if steps != 10 {
		failed = true
	}

	if failed {
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}
