// Topic: custom-errors | Concept: error-handling | Stage: 02-intermediate
// Run:   go run ./02-intermediate/error-handling/custom-errors
//
// Explore:
// * Any type with Error() string is an error
// * Custom error structs carry extra fields (code, field name, ...)
// * Use a pointer receiver and return *MyError, then match with errors.As
// * Add Unwrap() error so errors.Is/As can see the cause inside
// * Trap: returning a typed nil pointer as error gives a non-nil error
package main

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	Field string
	Value any
	Cause error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid %s = %v: %v", e.Field, e.Value, e.Cause)
}

func (e *ValidationError) Unwrap() error {
	return e.Cause
}

var ErrTooSmall = errors.New("too small")

func validateQty(qty int) error {
	if qty < 1 {
		return &ValidationError{Field: "qty", Value: qty, Cause: ErrTooSmall}
	}
	return nil
}

// The typed-nil trap
func badValidate() error {
	var e *ValidationError // nil pointer
	return e               // but the error interface is NOT nil
}

func main() {
	fmt.Println("Topic: custom errors")

	err := validateQty(0)
	fmt.Println("error:", err)

	var ve *ValidationError
	if errors.As(err, &ve) {
		fmt.Printf("errors.As -> field=%s value=%v\n", ve.Field, ve.Value)
	}
	fmt.Println("errors.Is(err, ErrTooSmall) via Unwrap:", errors.Is(err, ErrTooSmall))

	wrapped := fmt.Errorf("order 17: %w", err)
	fmt.Println("still found when wrapped:", errors.As(wrapped, &ve))

	fmt.Println("validateQty(3) == nil:", validateQty(3) == nil)
	fmt.Println("badValidate() == nil:", badValidate() == nil, "(typed nil trap)")
}
