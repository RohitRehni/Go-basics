// Topic: fibonacci-memo | Concept: dynamic-programming | Stage: 04-dsa
// Run:   go run ./04-dsa/dynamic-programming/fibonacci-memo
//
// Explore:
//   - Dynamic programming = overlapping subproblems + optimal substructure:
//     solve each subproblem once and reuse the answer
//   - Naive recursion recomputes the same values: O(2^n) calls
//   - Top-down memoization: recursion + a cache, O(n) time and space
//   - Bottom-up tabulation: fill a table from the smallest case up, O(n)
//   - Space optimisation: keep only the last two values, O(1) space
//   - Same pattern: climbing stairs, coin change, grid paths
package main

import (
	"fmt"
	"os"
	"time"
)

var calls int

func fibNaive(n int) int {
	calls++
	if n < 2 {
		return n
	}
	return fibNaive(n-1) + fibNaive(n-2)
}

func fibMemo(n int, memo map[int]int) int {
	calls++
	if n < 2 {
		return n
	}
	if v, ok := memo[n]; ok {
		return v
	}
	memo[n] = fibMemo(n-1, memo) + fibMemo(n-2, memo)
	return memo[n]
}

func fibTable(n int) int {
	if n < 2 {
		return n
	}
	table := make([]int, n+1)
	table[1] = 1
	for i := 2; i <= n; i++ {
		table[i] = table[i-1] + table[i-2]
	}
	return table[n]
}

func fibTwoVars(n int) int {
	a, b := 0, 1
	for range n {
		a, b = b, a+b
	}
	return a
}

// climbStairs: ways to climb n steps taking 1 or 2 at a time (Fibonacci in disguise).
func climbStairs(n int) int {
	one, two := 1, 1
	for i := 2; i <= n; i++ {
		one, two = one+two, one
	}
	return one
}

// coinChange returns the fewest coins that make amount, or -1.
func coinChange(coins []int, amount int) int {
	best := make([]int, amount+1)
	for a := 1; a <= amount; a++ {
		best[a] = -1
		for _, c := range coins {
			if c <= a && best[a-c] != -1 && (best[a] == -1 || best[a-c]+1 < best[a]) {
				best[a] = best[a-c] + 1
			}
		}
	}
	return best[amount]
}

func main() {
	fmt.Println("Topic: Fibonacci, memoization and tabulation")
	failed := false

	calls = 0
	start := time.Now()
	naive := fibNaive(30)
	fmt.Printf("naive  fib(30) = %d with %d calls in %v\n", naive, calls, time.Since(start).Round(time.Millisecond))

	calls = 0
	memo := fibMemo(30, map[int]int{})
	fmt.Printf("memo   fib(30) = %d with %d calls\n", memo, calls)
	fmt.Printf("table  fib(30) = %d\n", fibTable(30))
	fmt.Printf("2 vars fib(90) = %d\n", fibTwoVars(90))

	if naive != 832040 || memo != naive || fibTable(30) != naive || fibTwoVars(30) != naive {
		failed = true
	}
	if fibTwoVars(90) != 2880067194370816120 {
		failed = true
	}

	fmt.Println("ways to climb 10 stairs (1 or 2 steps):", climbStairs(10))
	if climbStairs(10) != 89 {
		failed = true
	}
	fmt.Println("fewest coins for 63 with [1 5 10 21 25]:", coinChange([]int{1, 5, 10, 21, 25}, 63))
	fmt.Println("fewest coins for 3 with [2]:", coinChange([]int{2}, 3))
	if coinChange([]int{1, 5, 10, 21, 25}, 63) != 3 || coinChange([]int{2}, 3) != -1 {
		failed = true
	}

	if failed {
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}
