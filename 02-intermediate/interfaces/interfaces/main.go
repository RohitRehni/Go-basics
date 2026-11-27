// Topic: interfaces | Concept: interfaces | Stage: 02-intermediate
// Run:   go run ./02-intermediate/interfaces/interfaces
//
// Explore:
// * An interface is a set of method signatures
// * Types satisfy interfaces implicitly: no "implements" keyword
// * A slice of an interface type can hold different concrete types
// * Pointer vs value receivers decide which type satisfies the interface
// * The empty interface `any` accepts every value
// * A nil interface vs an interface holding a nil pointer
// * Compile-time check: var _ Shape = Circle{}
package main

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Circle struct{ R float64 }
type Square struct{ Side float64 }

func (c Circle) Area() float64      { return math.Pi * c.R * c.R }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.R }
func (s Square) Area() float64      { return s.Side * s.Side }
func (s Square) Perimeter() float64 { return 4 * s.Side }

// Compile-time proof that both types satisfy Shape
var _ Shape = Circle{}
var _ Shape = Square{}

type Counter struct{ n int }

func (c *Counter) Area() float64      { c.n++; return 0 }
func (c *Counter) Perimeter() float64 { return 0 }

func total(shapes []Shape) float64 {
	sum := 0.0
	for _, s := range shapes {
		sum += s.Area()
	}
	return sum
}

func main() {
	fmt.Println("Topic: interfaces")

	shapes := []Shape{Circle{R: 1}, Square{Side: 2}}
	for _, s := range shapes {
		fmt.Printf("%-14T area = %6.2f, perimeter = %6.2f\n", s, s.Area(), s.Perimeter())
	}
	fmt.Printf("total area = %.2f\n", total(shapes))

	// Only *Counter has the methods, so only a pointer satisfies Shape
	var s Shape = &Counter{}
	s.Area()
	fmt.Printf("*Counter satisfies Shape: %T\n", s)
	// var bad Shape = Counter{} // compile error: Counter does not implement Shape (method Area has pointer receiver)

	var anything any = 42
	fmt.Printf("any holds %v of type %T\n", anything, anything)

	var nilIface Shape
	var nilPtr *Counter
	var holdsNilPtr Shape = nilPtr
	fmt.Println("nil interface == nil:", nilIface == nil)
	fmt.Println("interface holding nil pointer == nil:", holdsNilPtr == nil) // false: the type is set
}
