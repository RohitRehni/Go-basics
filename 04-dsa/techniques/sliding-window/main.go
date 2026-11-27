// Topic: sliding-window | Concept: techniques | Stage: 04-dsa
// Run:   go run ./04-dsa/techniques/sliding-window
//
// Explore:
// * A window [left, right] moves over the data; update the answer as items enter and leave
// * Fixed size k: add the new item, remove the one falling out: O(n) instead of O(n*k)
// * Variable size: grow right; shrink left while the window breaks the rule
// * Keep window state in a counter map or a running sum
// * Monotonic deque gives the max of every window in O(n)
// * Uses: moving averages, rate limiting, longest substring without repeats
package main

import (
	"fmt"
	"os"
	"slices"
)

// MaxSumWindow returns the largest sum of any k consecutive numbers.
func MaxSumWindow(nums []int, k int) int {
	sum := 0
	for i := range k {
		sum += nums[i]
	}
	best := sum
	for i := k; i < len(nums); i++ {
		sum += nums[i] - nums[i-k] // slide: one in, one out
		best = max(best, sum)
	}
	return best
}

// LongestUniqueSubstring returns the longest substring with no repeated character.
func LongestUniqueSubstring(s string) string {
	r := []rune(s)
	lastSeen := map[rune]int{}
	left, bestL, bestLen := 0, 0, 0
	for right, c := range r {
		if i, ok := lastSeen[c]; ok && i >= left {
			left = i + 1 // shrink past the earlier copy
		}
		lastSeen[c] = right
		if right-left+1 > bestLen {
			bestL, bestLen = left, right-left+1
		}
	}
	return string(r[bestL : bestL+bestLen])
}

// MinWindowAtLeast returns the shortest window length with sum >= target (0 if none).
func MinWindowAtLeast(nums []int, target int) int {
	best, sum, left := 0, 0, 0
	for right, v := range nums {
		sum += v
		for sum >= target {
			if best == 0 || right-left+1 < best {
				best = right - left + 1
			}
			sum -= nums[left]
			left++
		}
	}
	return best
}

// WindowMax returns the max of every window of size k with a monotonic deque of indexes.
func WindowMax(nums []int, k int) []int {
	var dq []int // indexes, values decreasing from front to back
	var out []int
	for i, v := range nums {
		for len(dq) > 0 && nums[dq[len(dq)-1]] <= v {
			dq = dq[:len(dq)-1] // smaller values can never be a max again
		}
		dq = append(dq, i)
		if dq[0] <= i-k {
			dq = dq[1:] // front fell out of the window
		}
		if i >= k-1 {
			out = append(out, nums[dq[0]])
		}
	}
	return out
}

func main() {
	fmt.Println("Topic: sliding window")
	failed := false

	nums := []int{2, 1, 5, 1, 3, 2, 9, 1}
	fmt.Println("max sum of 3 in a row in", nums, "=", MaxSumWindow(nums, 3))
	if MaxSumWindow(nums, 3) != 14 {
		failed = true
	}

	for _, tc := range [][2]string{{"abcabcbb", "abc"}, {"bbbbb", "b"}, {"pwwkew", "wke"}} {
		s, want := tc[0], tc[1]
		got := LongestUniqueSubstring(s)
		if got != want {
			failed = true
		}
		fmt.Printf("longest unique substring of %-10q = %q\n", s, got)
	}

	fmt.Println("shortest window with sum >= 7 in [2 3 1 2 4 3]:", MinWindowAtLeast([]int{2, 3, 1, 2, 4, 3}, 7))
	if MinWindowAtLeast([]int{2, 3, 1, 2, 4, 3}, 7) != 2 {
		failed = true
	}

	wm := WindowMax([]int{1, 3, -1, -3, 5, 3, 6, 7}, 3)
	fmt.Println("max of each window of 3:", wm)
	if !slices.Equal(wm, []int{3, 3, 5, 5, 6, 7}) {
		failed = true
	}

	if failed {
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}
