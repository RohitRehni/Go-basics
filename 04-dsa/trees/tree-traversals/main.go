// Topic: tree-traversals | Concept: trees | Stage: 04-dsa
// Run:   go run ./04-dsa/trees/tree-traversals
//
// Explore:
//   - Depth-first orders: pre-order (node, left, right), in-order (left, node, right),
//     post-order (left, right, node)
//   - Breadth-first / level-order: visit level by level with a queue
//   - Recursive versions are short; iterative versions use an explicit stack
//   - All traversals are O(n) time; DFS uses O(h) space, BFS uses O(width)
//   - Uses: pre-order copies a tree, post-order deletes or sizes it, level-order prints by depth
package main

import (
	"fmt"
	"os"
	"slices"
)

type Node struct {
	Val         string
	Left, Right *Node
}

func PreOrder(n *Node, visit func(string)) {
	if n == nil {
		return
	}
	visit(n.Val)
	PreOrder(n.Left, visit)
	PreOrder(n.Right, visit)
}

func InOrder(n *Node, visit func(string)) {
	if n == nil {
		return
	}
	InOrder(n.Left, visit)
	visit(n.Val)
	InOrder(n.Right, visit)
}

func PostOrder(n *Node, visit func(string)) {
	if n == nil {
		return
	}
	PostOrder(n.Left, visit)
	PostOrder(n.Right, visit)
	visit(n.Val)
}

// InOrderIterative uses an explicit stack instead of recursion.
func InOrderIterative(root *Node) []string {
	var out []string
	var stack []*Node
	cur := root
	for cur != nil || len(stack) > 0 {
		for cur != nil {
			stack = append(stack, cur)
			cur = cur.Left
		}
		cur = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		out = append(out, cur.Val)
		cur = cur.Right
	}
	return out
}

// LevelOrder returns the values grouped by depth.
func LevelOrder(root *Node) [][]string {
	if root == nil {
		return nil
	}
	var levels [][]string
	queue := []*Node{root}
	for len(queue) > 0 {
		var level []string
		next := []*Node{}
		for _, n := range queue {
			level = append(level, n.Val)
			if n.Left != nil {
				next = append(next, n.Left)
			}
			if n.Right != nil {
				next = append(next, n.Right)
			}
		}
		levels = append(levels, level)
		queue = next
	}
	return levels
}

func collect(walk func(*Node, func(string)), root *Node) []string {
	var out []string
	walk(root, func(v string) { out = append(out, v) })
	return out
}

func main() {
	fmt.Println("Topic: tree traversals")
	//        F
	//      /   \
	//     B     G
	//    / \     \
	//   A   D     I
	//      / \   /
	//     C   E H
	root := &Node{"F",
		&Node{"B", &Node{"A", nil, nil}, &Node{"D", &Node{"C", nil, nil}, &Node{"E", nil, nil}}},
		&Node{"G", nil, &Node{"I", &Node{"H", nil, nil}, nil}},
	}

	failed := false
	check := func(label string, got, want []string) {
		status := "ok"
		if !slices.Equal(got, want) {
			status = fmt.Sprintf("MISMATCH (want %v)", want)
			failed = true
		}
		fmt.Printf("  %-20s %v %s\n", label, got, status)
	}

	check("pre-order", collect(PreOrder, root), []string{"F", "B", "A", "D", "C", "E", "G", "I", "H"})
	check("in-order", collect(InOrder, root), []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"})
	check("in-order (stack)", InOrderIterative(root), []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"})
	check("post-order", collect(PostOrder, root), []string{"A", "C", "E", "D", "B", "H", "I", "G", "F"})

	levels := LevelOrder(root)
	fmt.Println("  level-order:")
	for depth, lvl := range levels {
		fmt.Printf("    depth %d: %v\n", depth, lvl)
	}
	if len(levels) != 4 || !slices.Equal(levels[3], []string{"C", "E", "H"}) {
		failed = true
	}

	if failed {
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}
