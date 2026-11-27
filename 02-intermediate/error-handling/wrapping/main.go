// Topic: wrapping | Concept: error-handling | Stage: 02-intermediate
// Run:   go run ./02-intermediate/error-handling/wrapping
//
// Explore:
// * fmt.Errorf("context: %w", err) wraps an error and keeps the original
// * %v instead of %w only copies the text: the chain is lost
// * errors.Is(err, target) walks the chain looking for a value
// * errors.As(err, &target) walks the chain looking for a type
// * errors.Unwrap peels one layer
// * errors.Join combines several errors into one
package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

var ErrNotFound = errors.New("not found")

func findUser(id int) error {
	if id != 1 {
		return ErrNotFound
	}
	return nil
}

func loadProfile(id int) error {
	if err := findUser(id); err != nil {
		return fmt.Errorf("loadProfile(%d): %w", id, err)
	}
	return nil
}

func handler(id int) error {
	if err := loadProfile(id); err != nil {
		return fmt.Errorf("handler: %w", err)
	}
	return nil
}

func main() {
	fmt.Println("Topic: wrapping errors")

	err := handler(7)
	fmt.Println("full message:", err)
	fmt.Println("errors.Is(err, ErrNotFound):", errors.Is(err, ErrNotFound))
	fmt.Println("err == ErrNotFound:", err == ErrNotFound) // false: it is wrapped
	fmt.Println("one layer unwrapped:", errors.Unwrap(err))

	lost := fmt.Errorf("with %%v: %v", ErrNotFound)
	fmt.Printf("errors.Is after wrapping with %%v: %v\n", errors.Is(lost, ErrNotFound))

	_, openErr := os.Open("/missing/config.yaml")
	wrapped := fmt.Errorf("reading config: %w", openErr)
	var pathErr *fs.PathError
	if errors.As(wrapped, &pathErr) {
		fmt.Printf("errors.As found *fs.PathError: op=%s path=%s\n", pathErr.Op, pathErr.Path)
	}
	fmt.Println("errors.Is(wrapped, fs.ErrNotExist):", errors.Is(wrapped, fs.ErrNotExist))

	joined := errors.Join(ErrNotFound, errors.New("timeout"))
	fmt.Printf("joined:\n%v\n", joined)
	fmt.Println("errors.Is(joined, ErrNotFound):", errors.Is(joined, ErrNotFound))
}
