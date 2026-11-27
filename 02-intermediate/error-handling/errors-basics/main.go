// Topic: errors-basics | Concept: error-handling | Stage: 02-intermediate
// Run:   go run ./02-intermediate/error-handling/errors-basics
//
// Explore:
// * error is a built-in interface: Error() string
// * Return the error as the last value; nil means success
// * Check right away: if err != nil { return ..., err }
// * errors.New and fmt.Errorf build simple errors
// * Sentinel errors: package-level vars you compare against
// * Library errors you meet every day: strconv.Atoi, os.Open
package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

var ErrDivideByZero = errors.New("divide by zero") // sentinel error

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, ErrDivideByZero
	}
	return a / b, nil
}

func parseAge(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("age %q is not a number", s)
	}
	if n < 0 {
		return 0, fmt.Errorf("age %d cannot be negative", n)
	}
	return n, nil
}

func main() {
	fmt.Println("Topic: error basics")

	if q, err := divide(10, 2); err == nil {
		fmt.Println("10 / 2 =", q)
	}

	_, err := divide(1, 0)
	if err != nil {
		fmt.Println("divide(1, 0) error:", err)
		fmt.Println("is the sentinel:", err == ErrDivideByZero)
	}

	for _, input := range []string{"42", "abc", "-5"} {
		age, err := parseAge(input)
		if err != nil {
			fmt.Printf("parseAge(%q) -> error: %v\n", input, err)
			continue
		}
		fmt.Printf("parseAge(%q) -> %d\n", input, age)
	}

	_, err = os.Open("/definitely/not/here.txt")
	fmt.Println("os.Open error:", err)

	// result := divide(1, 2) // compile error: assignment mismatch: 1 variable but divide returns 2 values
}
