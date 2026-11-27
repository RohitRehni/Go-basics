// Topic: define-structs | Concept: structs | Stage: 02-intermediate
// Run:   go run ./02-intermediate/structs/define-structs
//
// Explore:
// * Declaring a struct type and its fields
// * Struct literals: named fields vs positional
// * Zero values: an empty struct has every field at its zero value
// * Pointers to structs: p.Field works without (*p).Field
// * Anonymous structs for one-off data
// * Structs are values: assigning copies all fields
// * Comparing structs with == (only if every field is comparable)
package main

import "fmt"

type Employee struct {
	ID     int
	Name   string
	Salary float64
	Active bool
}

func main() {
	fmt.Println("Topic: define structs")

	// Named fields: order does not matter, missing fields get zero values
	e1 := Employee{Name: "Asha", ID: 1, Salary: 50000}
	fmt.Printf("e1 = %+v\n", e1)

	// Positional: every field, in order (fragile if the struct changes)
	e2 := Employee{2, "Ravi", 42000, true}
	fmt.Printf("e2 = %+v\n", e2)

	// Zero value
	var empty Employee
	fmt.Printf("zero value = %+v (Name is %q)\n", empty, empty.Name)

	// Pointer to a struct: Go dereferences automatically
	p := &e1
	p.Salary += 5000
	fmt.Println("after raise via pointer, e1.Salary =", e1.Salary)

	// Assignment copies
	copyOf := e2
	copyOf.Name = "Changed"
	fmt.Printf("e2.Name = %s, copyOf.Name = %s\n", e2.Name, copyOf.Name)

	// Comparison
	fmt.Println("e2 == Employee{2, \"Ravi\", 42000, true}:", e2 == Employee{2, "Ravi", 42000, true})

	// Anonymous struct
	point := struct{ X, Y int }{X: 3, Y: 4}
	fmt.Printf("anonymous struct = %+v\n", point)

	// e3 := Employee{1, "x"} // compile error: too few values in struct literal
}
