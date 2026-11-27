// Topic: hash-map | Concept: hash-tables | Stage: 04-dsa
// Run:   go run ./04-dsa/hash-tables/hash-map
//
// Explore:
//   - A hash function turns a key into a bucket index: hash(key) % len(buckets)
//   - Collisions: two keys land in the same bucket. Separate chaining keeps a
//     small list per bucket; open addressing probes for the next free slot
//   - Load factor = items / buckets. Above ~0.75, resize (double) and rehash everything
//   - Get/Put/Delete: O(1) average, O(n) worst case if every key collides
//   - Resizing costs O(n) once, but doubling makes it O(1) amortized per Put
//   - Go's built-in map is a hash table (Swiss tables since Go 1.24); iteration order is random
package main

import (
	"fmt"
	"hash/maphash"
	"math/rand/v2"
	"os"
)

type entry[K comparable, V any] struct {
	key K
	val V
}

type HashMap[K comparable, V any] struct {
	buckets [][]entry[K, V]
	size    int
	seed    maphash.Seed
	resizes int
}

func NewHashMap[K comparable, V any]() *HashMap[K, V] {
	return &HashMap[K, V]{buckets: make([][]entry[K, V], 8), seed: maphash.MakeSeed()}
}

func (m *HashMap[K, V]) index(key K, n int) int {
	return int(maphash.Comparable(m.seed, key) % uint64(n))
}

func (m *HashMap[K, V]) Put(key K, val V) {
	i := m.index(key, len(m.buckets))
	for j := range m.buckets[i] {
		if m.buckets[i][j].key == key { // update existing key
			m.buckets[i][j].val = val
			return
		}
	}
	m.buckets[i] = append(m.buckets[i], entry[K, V]{key, val})
	m.size++
	if float64(m.size)/float64(len(m.buckets)) > 0.75 {
		m.resize()
	}
}

func (m *HashMap[K, V]) Get(key K) (V, bool) {
	for _, e := range m.buckets[m.index(key, len(m.buckets))] {
		if e.key == key {
			return e.val, true
		}
	}
	var zero V
	return zero, false
}

func (m *HashMap[K, V]) Delete(key K) bool {
	i := m.index(key, len(m.buckets))
	b := m.buckets[i]
	for j, e := range b {
		if e.key == key {
			b[j] = b[len(b)-1] // swap with last, then shrink
			m.buckets[i] = b[:len(b)-1]
			m.size--
			return true
		}
	}
	return false
}

// resize doubles the bucket count and rehashes every entry.
func (m *HashMap[K, V]) resize() {
	nb := make([][]entry[K, V], 2*len(m.buckets))
	for _, b := range m.buckets {
		for _, e := range b {
			i := m.index(e.key, len(nb))
			nb[i] = append(nb[i], e)
		}
	}
	m.buckets = nb
	m.resizes++
}

func (m *HashMap[K, V]) longestChain() int {
	longest := 0
	for _, b := range m.buckets {
		longest = max(longest, len(b))
	}
	return longest
}

func main() {
	fmt.Println("Topic: hash map with separate chaining")
	failed := false

	m := NewHashMap[string, int]()
	m.Put("apple", 3)
	m.Put("banana", 5)
	m.Put("apple", 7) // update
	v, ok := m.Get("apple")
	fmt.Println("Get(apple):", v, ok, "| size:", m.size)
	_, ok = m.Get("cherry")
	fmt.Println("Get(cherry) found:", ok)
	fmt.Println("Delete(banana):", m.Delete("banana"), "| size:", m.size)
	if v != 7 || m.size != 1 {
		failed = true
	}

	// Stress test against Go's built-in map
	ours := NewHashMap[int, int]()
	builtin := map[int]int{}
	r := rand.New(rand.NewPCG(3, 4))
	for range 20_000 {
		k := r.IntN(5_000)
		switch r.IntN(3) {
		case 0, 1:
			ours.Put(k, k*10)
			builtin[k] = k * 10
		case 2:
			ours.Delete(k)
			delete(builtin, k)
		}
	}
	mismatches := 0
	for k := range 5_000 {
		a, okA := ours.Get(k)
		b, okB := builtin[k]
		if a != b || okA != okB {
			mismatches++
		}
	}
	fmt.Printf("stress test: size %d (built-in %d), buckets %d, resizes %d, longest chain %d, mismatches %d\n",
		ours.size, len(builtin), len(ours.buckets), ours.resizes, ours.longestChain(), mismatches)
	if mismatches != 0 || ours.size != len(builtin) {
		failed = true
	}

	if failed {
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}
