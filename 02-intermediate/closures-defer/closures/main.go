// Topic: closures | Concept: closures-defer | Stage: 02-intermediate
// Run:   go run ./02-intermediate/closures-defer/closures
//
// Explore:
// * Functions are values: assign them, pass them, return them
// * A closure captures variables from its surrounding scope by reference
// * Each call of a closure factory gets its own captured state
// * Function types: type Op func(int, int) int
// * Since Go 1.22 each loop iteration has a fresh loop variable
// * Closures for middleware-style wrapping
package main

import (
	"fmt"
	"strings"
)

type Op func(int, int) int

func counter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

func apply(a, b int, op Op) int {
	return op(a, b)
}

func withLogging(name string, f func(string) string) func(string) string {
	return func(in string) string {
		out := f(in)
		fmt.Printf("  [log] %s(%q) = %q\n", name, in, out)
		return out
	}
}

func main() {
	fmt.Println("Topic: closures")

	c1 := counter()
	c2 := counter()
	fmt.Println("c1:", c1(), c1(), c1())
	fmt.Println("c2 has its own count:", c2())

	add := func(a, b int) int { return a + b }
	fmt.Println("apply(3, 4, add):", apply(3, 4, add))
	fmt.Println("apply with inline func:", apply(3, 4, func(a, b int) int { return a * b }))

	total := 0
	addTo := func(n int) { total += n } // captures total
	addTo(5)
	addTo(10)
	fmt.Println("total changed by closure:", total)

	var funcs []func() int
	for i := 0; i < 3; i++ {
		funcs = append(funcs, func() int { return i })
	}
	for _, f := range funcs {
		fmt.Print(f(), " ") // 0 1 2 since Go 1.22 (was 3 3 3 before)
	}
	fmt.Println()

	shout := withLogging("shout", strings.ToUpper)
	shout("hello")
}
