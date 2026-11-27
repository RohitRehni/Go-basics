// Topic: stack | Concept: stacks-queues | Stage: 04-dsa
// Run:   go run ./04-dsa/stacks-queues/stack
//
// Explore:
// * LIFO: last in, first out. Push, Pop and Peek are all O(1)
// * A slice is the natural backing store: append to push, shrink to pop
// * Pop on an empty stack must be handled: return (value, ok)
// * Uses: balanced brackets, undo history, expression evaluation, DFS, the call stack
// * Postfix (RPN) evaluation: push numbers, pop two for each operator
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) }

func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	v := s.items[len(s.items)-1]
	s.items[len(s.items)-1] = zero // let the GC free it
	s.items = s.items[:len(s.items)-1]
	return v, true
}

func (s *Stack[T]) Peek() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

func (s *Stack[T]) Len() int { return len(s.items) }

// Balanced reports whether every bracket is closed in the right order.
func Balanced(text string) bool {
	pairs := map[rune]rune{')': '(', ']': '[', '}': '{'}
	var st Stack[rune]
	for _, r := range text {
		switch r {
		case '(', '[', '{':
			st.Push(r)
		case ')', ']', '}':
			top, ok := st.Pop()
			if !ok || top != pairs[r] {
				return false
			}
		}
	}
	return st.Len() == 0
}

// EvalRPN evaluates a postfix expression such as "3 4 + 2 *".
func EvalRPN(expr string) (int, error) {
	var st Stack[int]
	for _, tok := range strings.Fields(expr) {
		switch tok {
		case "+", "-", "*", "/":
			b, ok1 := st.Pop()
			a, ok2 := st.Pop()
			if !ok1 || !ok2 {
				return 0, fmt.Errorf("not enough operands for %q", tok)
			}
			switch tok {
			case "+":
				st.Push(a + b)
			case "-":
				st.Push(a - b)
			case "*":
				st.Push(a * b)
			case "/":
				if b == 0 {
					return 0, fmt.Errorf("division by zero")
				}
				st.Push(a / b)
			}
		default:
			n, err := strconv.Atoi(tok)
			if err != nil {
				return 0, fmt.Errorf("bad token %q", tok)
			}
			st.Push(n)
		}
	}
	if st.Len() != 1 {
		return 0, fmt.Errorf("expression left %d values", st.Len())
	}
	v, _ := st.Pop()
	return v, nil
}

func main() {
	fmt.Println("Topic: stack")
	failed := false

	var s Stack[string]
	s.Push("first")
	s.Push("second")
	s.Push("third")
	top, _ := s.Peek()
	fmt.Println("peek:", top, "| len:", s.Len())
	for s.Len() > 0 {
		v, _ := s.Pop()
		fmt.Println("  pop:", v)
	}
	_, ok := s.Pop()
	fmt.Println("pop on empty -> ok =", ok)

	fmt.Println("balanced brackets:")
	for _, tc := range []struct {
		text string
		want bool
	}{
		{"([]{})", true},
		{"func() { x[0] }", true},
		{"([)]", false},
		{"((", false},
		{"", true},
		{"}", false},
	} {
		text, want := tc.text, tc.want
		got := Balanced(text)
		status := "ok"
		if got != want {
			status = "MISMATCH"
			failed = true
		}
		fmt.Printf("  %-18q -> %-5v %s\n", text, got, status)
	}

	fmt.Println("postfix (RPN) evaluation:")
	for _, tc := range []struct {
		expr string
		want int
	}{{"3 4 + 2 *", 14}, {"5 1 2 + 4 * + 3 -", 14}, {"10 2 /", 5}} {
		expr, want := tc.expr, tc.want
		got, err := EvalRPN(expr)
		status := "ok"
		if err != nil || got != want {
			status = "MISMATCH"
			failed = true
		}
		fmt.Printf("  %-20q = %-3d %s\n", expr, got, status)
	}
	if _, err := EvalRPN("1 +"); err != nil {
		fmt.Println("  \"1 +\" error:", err)
	} else {
		failed = true
	}

	if failed {
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}
