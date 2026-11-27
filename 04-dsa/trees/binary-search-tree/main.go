// Topic: binary-search-tree | Concept: trees | Stage: 04-dsa
// Run:   go run ./04-dsa/trees/binary-search-tree
//
// Explore:
// * BST rule: every key in the left subtree < node < every key in the right subtree
// * Insert / search / delete: O(h) where h = height
// * Balanced tree h = O(log n); inserting sorted keys makes a "linked list", h = O(n)
// * Delete has 3 cases: leaf, one child, two children (replace with the in-order successor)
// * In-order traversal visits keys in sorted order
// * Self-balancing trees (AVL, red-black) keep h = O(log n) automatically
package main

import (
	"cmp"
	"fmt"
	"math/rand/v2"
	"os"
	"slices"
)

type bstNode[K cmp.Ordered] struct {
	key         K
	left, right *bstNode[K]
}

type BST[K cmp.Ordered] struct {
	root *bstNode[K]
	size int
}

func (t *BST[K]) Insert(key K) bool {
	link := &t.root
	for *link != nil {
		switch {
		case key < (*link).key:
			link = &(*link).left
		case key > (*link).key:
			link = &(*link).right
		default:
			return false // duplicate
		}
	}
	*link = &bstNode[K]{key: key}
	t.size++
	return true
}

func (t *BST[K]) Contains(key K) bool {
	n := t.root
	for n != nil {
		switch {
		case key < n.key:
			n = n.left
		case key > n.key:
			n = n.right
		default:
			return true
		}
	}
	return false
}

func (t *BST[K]) Delete(key K) bool {
	var deleted bool
	t.root, deleted = deleteNode(t.root, key)
	if deleted {
		t.size--
	}
	return deleted
}

func deleteNode[K cmp.Ordered](n *bstNode[K], key K) (*bstNode[K], bool) {
	if n == nil {
		return nil, false
	}
	var deleted bool
	switch {
	case key < n.key:
		n.left, deleted = deleteNode(n.left, key)
		return n, deleted
	case key > n.key:
		n.right, deleted = deleteNode(n.right, key)
		return n, deleted
	}
	// Found it.
	if n.left == nil {
		return n.right, true // leaf or only a right child
	}
	if n.right == nil {
		return n.left, true // only a left child
	}
	succ := n.right // two children: smallest key in the right subtree
	for succ.left != nil {
		succ = succ.left
	}
	n.key = succ.key
	n.right, _ = deleteNode(n.right, succ.key)
	return n, true
}

func (t *BST[K]) InOrder() []K {
	var out []K
	var walk func(*bstNode[K])
	walk = func(n *bstNode[K]) {
		if n == nil {
			return
		}
		walk(n.left)
		out = append(out, n.key)
		walk(n.right)
	}
	walk(t.root)
	return out
}

func height[K cmp.Ordered](n *bstNode[K]) int {
	if n == nil {
		return 0
	}
	return 1 + max(height(n.left), height(n.right))
}

func (t *BST[K]) Min() (K, bool) {
	var zero K
	if t.root == nil {
		return zero, false
	}
	n := t.root
	for n.left != nil {
		n = n.left
	}
	return n.key, true
}

func main() {
	fmt.Println("Topic: binary search tree")
	failed := false

	var t BST[int]
	for _, k := range []int{50, 30, 70, 20, 40, 60, 80, 35} {
		t.Insert(k)
	}
	fmt.Println("in-order:", t.InOrder(), "| height:", height(t.root))
	fmt.Println("contains 40:", t.Contains(40), "| contains 45:", t.Contains(45))
	fmt.Println("insert 40 again (duplicate):", t.Insert(40))
	minKey, _ := t.Min()
	fmt.Println("min:", minKey)

	t.Delete(20) // leaf
	t.Delete(30) // two children
	t.Delete(70) // two children
	fmt.Println("after deleting 20, 30, 70:", t.InOrder())
	if !slices.Equal(t.InOrder(), []int{35, 40, 50, 60, 80}) {
		failed = true
	}

	// Balanced vs degenerate height
	var random, sorted BST[int]
	r := rand.New(rand.NewPCG(5, 6))
	for _, k := range r.Perm(1023) {
		random.Insert(k)
	}
	for k := range 1023 {
		sorted.Insert(k)
	}
	fmt.Printf("height with 1023 random keys: %d | with 1023 sorted keys: %d (a linked list!)\n",
		height(random.root), height(sorted.root))

	// Random operations vs a sorted-slice model
	var tree BST[int]
	model := map[int]bool{}
	for range 5_000 {
		k := r.IntN(500)
		if r.IntN(3) == 0 {
			tree.Delete(k)
			delete(model, k)
		} else {
			tree.Insert(k)
			model[k] = true
		}
	}
	want := make([]int, 0, len(model))
	for k := range model {
		want = append(want, k)
	}
	slices.Sort(want)
	ok := slices.Equal(tree.InOrder(), want) && tree.size == len(want)
	fmt.Println("5000 random inserts/deletes match a set model:", ok)
	if !ok {
		failed = true
	}

	if failed {
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}
