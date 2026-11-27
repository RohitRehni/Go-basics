// Topic: embedding | Concept: structs | Stage: 02-intermediate
// Run:   go run ./02-intermediate/structs/embedding
//
// Explore:
// * Go has no inheritance; it uses composition by embedding
// * Embedded fields and methods are "promoted": m.Name instead of m.Person.Name
// * The outer type can shadow a promoted method with its own
// * The embedded value is still reachable by its type name: m.Person
// * Embedding is not subtyping: a Manager is not a Person
package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

func (p Person) Greet() string {
	return "Hi, I am " + p.Name
}

func (p Person) Describe() string {
	return fmt.Sprintf("%s (%d)", p.Name, p.Age)
}

type Manager struct {
	Person // embedded: no field name
	Team   []string
}

// Shadows Person.Greet
func (m Manager) Greet() string {
	return m.Person.Greet() + " and I manage " + fmt.Sprint(len(m.Team)) + " people"
}

func describe(p Person) string {
	return p.Describe()
}

func main() {
	fmt.Println("Topic: embedding")

	m := Manager{
		Person: Person{Name: "Meera", Age: 40},
		Team:   []string{"Asha", "Ravi"},
	}

	fmt.Println("promoted field m.Name:", m.Name)
	fmt.Println("promoted method m.Describe():", m.Describe())
	fmt.Println("shadowed method m.Greet():", m.Greet())
	fmt.Println("original method m.Person.Greet():", m.Person.Greet())

	fmt.Println("pass the embedded value:", describe(m.Person))
	// describe(m) // compile error: cannot use m (Manager) as Person value
}
