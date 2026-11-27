// Topic: two-pointers | Concept: techniques | Stage: 04-dsa
// Run:   go run ./04-dsa/techniques/two-pointers
//
// Explore:
// * Two indexes moving through the data replace a nested O(n^2) loop with O(n)
// * Opposite ends (left++ / right--): pair sums in sorted data, palindromes, container with most water
// * Same direction (slow / fast): remove duplicates in place, partition, merge sorted lists
// * Usually needs sorted input, or a property that tells you which pointer to move
// * Space O(1): works in place
package main

import (
	"fmt"
	"os"
	"slices"
	"unicode"
)

// PairSum finds two numbers in sorted nums that add up to target.
func PairSum(nums []int, target int) (int, int, bool) {
	l, r := 0, len(nums)-1
	for l < r {
		sum := nums[l] + nums[r]
		switch {
		case sum == target:
			return nums[l], nums[r], true
		case sum < target:
			l++ // need a bigger sum
		default:
			r-- // need a smaller sum
		}
	}
	return 0, 0, false
}

// IsPalindrome ignores case and anything that is not a letter or digit.
func IsPalindrome(s string) bool {
	r := []rune(s)
	l, h := 0, len(r)-1
	for l < h {
		if !unicode.IsLetter(r[l]) && !unicode.IsDigit(r[l]) {
			l++
			continue
		}
		if !unicode.IsLetter(r[h]) && !unicode.IsDigit(r[h]) {
			h--
			continue
		}
		if unicode.ToLower(r[l]) != unicode.ToLower(r[h]) {
			return false
		}
		l++
		h--
	}
	return true
}

// Dedupe removes duplicates from a sorted slice in place and returns the new length.
func Dedupe(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	slow := 0
	for fast := 1; fast < len(nums); fast++ {
		if nums[fast] != nums[slow] {
			slow++
			nums[slow] = nums[fast]
		}
	}
	return slow + 1
}

// MaxWater: two vertical lines and the x-axis form a container; find the largest area.
func MaxWater(heights []int) int {
	l, r := 0, len(heights)-1
	best := 0
	for l < r {
		best = max(best, (r-l)*min(heights[l], heights[r]))
		if heights[l] < heights[r] { // the shorter line limits the area: move it
			l++
		} else {
			r--
		}
	}
	return best
}

func main() {
	fmt.Println("Topic: two pointers")
	failed := false

	sorted := []int{1, 3, 4, 6, 8, 11, 15}
	a, b, ok := PairSum(sorted, 14)
	fmt.Printf("pair in %v summing to 14: %d + %d (%v)\n", sorted, a, b, ok)
	_, _, ok2 := PairSum(sorted, 100)
	fmt.Println("pair summing to 100 found:", ok2)
	if !ok || a+b != 14 || ok2 {
		failed = true
	}

	for _, tc := range []struct {
		s    string
		want bool
	}{{"A man, a plan, a canal: Panama", true}, {"race a car", false}, {"No 'x' in Nixon", true}} {
		s, want := tc.s, tc.want
		got := IsPalindrome(s)
		if got != want {
			failed = true
		}
		fmt.Printf("palindrome %-34q %v\n", s, got)
	}

	dup := []int{1, 1, 2, 3, 3, 3, 4, 5, 5}
	n := Dedupe(dup)
	fmt.Println("dedupe in place:", dup[:n])
	if !slices.Equal(dup[:n], []int{1, 2, 3, 4, 5}) {
		failed = true
	}

	h := []int{1, 8, 6, 2, 5, 4, 8, 3, 7}
	fmt.Println("max water for", h, "=", MaxWater(h))
	if MaxWater(h) != 49 {
		failed = true
	}

	if failed {
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}
