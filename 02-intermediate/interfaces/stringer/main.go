// Topic: stringer | Concept: interfaces | Stage: 02-intermediate
// Run:   go run ./02-intermediate/interfaces/stringer
//
// Explore:
// * fmt.Stringer: any type with String() string prints itself that way
// * Implementing String on an enum-like type built with iota
// * Value vs pointer receiver for String decides when fmt uses it
// * The error interface is the same idea: Error() string
// * Beware: calling fmt.Sprint(t) inside t.String() recurses forever
package main

import "fmt"

type Status int

const (
	Pending Status = iota
	Active
	Closed
)

func (s Status) String() string {
	switch s {
	case Pending:
		return "Pending"
	case Active:
		return "Active"
	case Closed:
		return "Closed"
	default:
		return fmt.Sprintf("Status(%d)", int(s))
	}
}

type Money struct {
	Cents int64
}

func (m Money) String() string {
	return fmt.Sprintf("$%d.%02d", m.Cents/100, m.Cents%100)
}

type Order struct {
	ID     int
	Total  Money
	Status Status
}

func main() {
	fmt.Println("Topic: fmt.Stringer")

	fmt.Println("Active prints as:", Active)
	fmt.Println("unknown value:", Status(7))
	fmt.Printf("%%v = %v, %%d = %d\n", Closed, Closed)

	o := Order{ID: 101, Total: Money{Cents: 123456}, Status: Active}
	fmt.Printf("order: %+v\n", o)

	var s fmt.Stringer = Money{Cents: 99}
	fmt.Println("as fmt.Stringer:", s.String())
}
