// Topic: type-switch | Concept: interfaces | Stage: 02-intermediate
// Run:   go run ./02-intermediate/interfaces/type-switch
//
// Explore:
// * Type assertion: v.(T) panics if wrong; v, ok := v.(T) does not
// * Type switch: switch x := v.(type) { case int: ... }
// * Several types in one case keep x as the interface type
// * Asserting to another interface: does this value also have method M?
package main

import (
	"fmt"
	"strings"
)

type Speaker interface {
	Speak() string
}

type Dog struct{}

func (Dog) Speak() string { return "Woof" }

func describe(v any) string {
	switch x := v.(type) {
	case nil:
		return "nil"
	case int:
		return fmt.Sprintf("int, doubled = %d", x*2)
	case string:
		return fmt.Sprintf("string, upper = %s", strings.ToUpper(x))
	case bool, float64:
		return fmt.Sprintf("bool or float64 (x is still any): %v", x)
	case Speaker:
		return "a Speaker that says " + x.Speak()
	case []int:
		return fmt.Sprintf("slice of %d ints", len(x))
	default:
		return fmt.Sprintf("something else: %T", x)
	}
}

func main() {
	fmt.Println("Topic: type assertions and type switches")

	var v any = "hello"

	s := v.(string)
	fmt.Println("v.(string) =", s)

	n, ok := v.(int)
	fmt.Printf("v.(int) with ok: n = %d, ok = %v\n", n, ok)
	// _ = v.(int) // runtime panic: interface conversion: interface {} is string, not int

	for _, value := range []any{nil, 21, "go", true, 3.5, Dog{}, []int{1, 2, 3}, map[string]int{}} {
		fmt.Printf("%-22v -> %s\n", value, describe(value))
	}

	if sp, ok := any(Dog{}).(Speaker); ok {
		fmt.Println("Dog also satisfies Speaker:", sp.Speak())
	}
}
