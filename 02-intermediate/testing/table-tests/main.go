// Topic: table-tests | Concept: testing | Stage: 02-intermediate
// Run:   go run ./02-intermediate/testing/table-tests
// Test:  go test ./02-intermediate/testing/table-tests
// Each:  go test -v ./02-intermediate/testing/table-tests   (show each case)
//
// Explore:
// * Tests live in *_test.go files next to the code, same package
// * func TestXxx(t *testing.T); t.Errorf reports and continues, t.Fatalf stops
// * Table-driven tests: a slice of cases, one loop, t.Run per case
// * go test -run 'TestIsPalindrome/mixed' runs one sub-test
// * Test files are ignored by go run and go build
package main

import (
	"fmt"
	"strings"
	"unicode"
)

// IsPalindrome ignores case, spaces and punctuation.
func IsPalindrome(s string) bool {
	var letters []rune
	for _, r := range strings.ToLower(s) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			letters = append(letters, r)
		}
	}
	for i, j := 0, len(letters)-1; i < j; i, j = i+1, j-1 {
		if letters[i] != letters[j] {
			return false
		}
	}
	return true
}

func main() {
	fmt.Println("Topic: table-driven tests")
	for _, s := range []string{"racecar", "A man, a plan, a canal: Panama", "golang"} {
		fmt.Printf("IsPalindrome(%q) = %v\n", s, IsPalindrome(s))
	}
	fmt.Println("Now run: go test -v ./02-intermediate/testing/table-tests")
}
