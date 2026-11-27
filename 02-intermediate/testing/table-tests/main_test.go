package main

import "testing"

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "empty", input: "", want: true},
		{name: "single letter", input: "a", want: true},
		{name: "simple", input: "racecar", want: true},
		{name: "mixed case and punctuation", input: "A man, a plan, a canal: Panama", want: true},
		{name: "digits", input: "12321", want: true},
		{name: "not a palindrome", input: "golang", want: false},
		{name: "unicode", input: "नमन", want: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := IsPalindrome(tc.input)
			if got != tc.want {
				t.Errorf("IsPalindrome(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}
