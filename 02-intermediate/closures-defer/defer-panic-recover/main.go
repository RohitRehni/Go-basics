// Topic: defer-panic-recover | Concept: closures-defer | Stage: 02-intermediate
// Run:   go run ./02-intermediate/closures-defer/defer-panic-recover
//
// Explore:
// * defer runs a call when the surrounding function returns
// * Several defers run in LIFO (stack) order
// * Deferred arguments are evaluated immediately, not at run time
// * Deferred closures can change named return values
// * panic stops normal flow and runs the defers on the way up
// * recover() only works inside a deferred function
// * Use errors for expected failures; panic for programmer bugs
package main

import (
	"errors"
	"fmt"
)

func order() {
	for i := 1; i <= 3; i++ {
		defer fmt.Println("  deferred", i)
	}
	fmt.Println("  body done")
}

func argsEvaluatedNow() {
	x := 1
	defer fmt.Println("  deferred sees x =", x) // 1, captured now
	x = 99
	fmt.Println("  body set x =", x)
}

func doubleOnReturn() (result int) {
	defer func() { result *= 2 }()
	return 21
}

func safeDivide(a, b int) (q int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered: %v", r)
		}
	}()
	return a / b, nil // panics when b == 0
}

func mustPositive(n int) int {
	if n < 0 {
		panic(errors.New("negative value"))
	}
	return n
}

func main() {
	fmt.Println("Topic: defer, panic and recover")

	fmt.Println("LIFO order:")
	order()

	fmt.Println("argument evaluation:")
	argsEvaluatedNow()

	fmt.Println("named result changed by defer:", doubleOnReturn())

	q, err := safeDivide(10, 2)
	fmt.Println("safeDivide(10, 2):", q, err)
	q, err = safeDivide(1, 0)
	fmt.Println("safeDivide(1, 0):", q, err)

	func() {
		defer func() {
			fmt.Println("recovered from mustPositive:", recover())
		}()
		mustPositive(-1)
		fmt.Println("never printed")
	}()

	fmt.Println("program continues normally")
}
