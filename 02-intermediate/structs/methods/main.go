// Topic: methods | Concept: structs | Stage: 02-intermediate
// Run:   go run ./02-intermediate/structs/methods
//
// Explore:
// * A method is a function with a receiver: func (r Rect) Area() float64
// * Value receiver: works on a copy, cannot change the caller's struct
// * Pointer receiver: can change the struct and avoids copying big structs
// * Go calls pointer methods on addressable values for you (c.Inc() on a var)
// * Methods can be declared on any named type you define, not only structs
// * Rule of thumb: if one method needs a pointer receiver, use pointers for all
package main

import "fmt"

type Rect struct {
	Width, Height float64
}

// Value receiver: reads only
func (r Rect) Area() float64 {
	return r.Width * r.Height
}

// Value receiver that tries to modify: changes only the copy
func (r Rect) DoubleCopy() {
	r.Width *= 2
}

// Pointer receiver: modifies the original
func (r *Rect) Double() {
	r.Width *= 2
	r.Height *= 2
}

// Methods on a non-struct named type
type Celsius float64

func (c Celsius) Fahrenheit() float64 {
	return float64(c)*9/5 + 32
}

func main() {
	fmt.Println("Topic: methods")

	r := Rect{Width: 3, Height: 4}
	fmt.Println("area =", r.Area())

	r.DoubleCopy()
	fmt.Printf("after DoubleCopy (value receiver): %+v\n", r)

	r.Double() // Go takes &r for you
	fmt.Printf("after Double (pointer receiver):   %+v, area = %v\n", r, r.Area())

	rp := &Rect{Width: 1, Height: 1}
	fmt.Println("value method on a pointer works too:", rp.Area())

	var temp Celsius = 100
	fmt.Printf("%v°C = %v°F\n", temp, temp.Fahrenheit())

	// Rect{1, 2}.Double() // compile error: cannot call pointer method Double on Rect (not addressable)
}
