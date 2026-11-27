// Topic: constraints | Concept: generics | Stage: 02-intermediate
// Run:   go run ./02-intermediate/generics/constraints
//
// Explore:
// * A constraint is an interface that limits which types T may be
// * Type sets: interface{ int | float64 }
// * ~int also accepts named types whose underlying type is int
// * cmp.Ordered from the standard library: anything with < >
// * Constraints mixing a method and a type set
// * Generic types with a constraint: type Pair[K comparable, V any]
package main

import (
	"cmp"
	"fmt"
	"strconv"
)

type Number interface {
	~int | ~int64 | ~float64
}

func Sum[T Number](values []T) T {
	var total T
	for _, v := range values {
		total += v
	}
	return total
}

func Max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// A named type with underlying int: allowed by ~int
type Rupees int

// Constraint with a method AND a type set
type StringableInt interface {
	~int
	String() string
}

type Code int

func (c Code) String() string { return "CODE-" + strconv.Itoa(int(c)) }

func Describe[T StringableInt](v T) string {
	return fmt.Sprintf("%s (raw %d)", v.String(), int(v))
}

type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

func main() {
	fmt.Println("Topic: generic constraints")

	fmt.Println("Sum ints:", Sum([]int{1, 2, 3}))
	fmt.Println("Sum floats:", Sum([]float64{1.5, 2.25}))
	fmt.Println("Sum Rupees (~int):", Sum([]Rupees{100, 250}))

	fmt.Println("Max(3, 7):", Max(3, 7))
	fmt.Println("Max(\"apple\", \"banana\"):", Max("apple", "banana"))

	fmt.Println("Describe(Code(42)):", Describe(Code(42)))

	p := Pair[string, int]{Key: "age", Value: 30}
	fmt.Printf("pair: %+v\n", p)

	// Sum([]string{"a"}) // compile error: string does not satisfy Number
}
