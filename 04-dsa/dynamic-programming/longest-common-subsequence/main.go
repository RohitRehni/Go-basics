// Topic: longest-common-subsequence | Concept: dynamic-programming | Stage: 04-dsa
// Run:   go run ./04-dsa/dynamic-programming/longest-common-subsequence
//
// Explore:
//   - A subsequence keeps order but may skip characters ("ace" is in "abcde")
//   - dp[i][j] = LCS length of a[:i] and b[:j]:
//     equal chars -> dp[i-1][j-1] + 1, else max(dp[i-1][j], dp[i][j-1])
//   - Time and space O(m * n); walk back through the table to rebuild the string
//   - Edit (Levenshtein) distance uses the same table shape: insert, delete, replace
//   - Real uses: diff tools, DNA alignment, spell checkers
//   - Work on []rune, not bytes, so non-ASCII text works
package main

import (
	"fmt"
	"os"
	"slices"
)

func LCS(a, b string) string {
	x, y := []rune(a), []rune(b)
	dp := make([][]int, len(x)+1)
	for i := range dp {
		dp[i] = make([]int, len(y)+1)
	}
	for i := 1; i <= len(x); i++ {
		for j := 1; j <= len(y); j++ {
			if x[i-1] == y[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}
	var out []rune // walk back from the bottom-right corner
	for i, j := len(x), len(y); i > 0 && j > 0; {
		switch {
		case x[i-1] == y[j-1]:
			out = append(out, x[i-1])
			i--
			j--
		case dp[i-1][j] >= dp[i][j-1]:
			i--
		default:
			j--
		}
	}
	slices.Reverse(out)
	return string(out)
}

func EditDistance(a, b string) int {
	x, y := []rune(a), []rune(b)
	prev := make([]int, len(y)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(x); i++ {
		cur := make([]int, len(y)+1)
		cur[0] = i
		for j := 1; j <= len(y); j++ {
			cost := 1
			if x[i-1] == y[j-1] {
				cost = 0
			}
			cur[j] = min(prev[j]+1, cur[j-1]+1, prev[j-1]+cost) // delete, insert, replace
		}
		prev = cur
	}
	return prev[len(y)]
}

func isSubsequence(sub, s string) bool {
	r := []rune(sub)
	i := 0
	for _, c := range s {
		if i < len(r) && r[i] == c {
			i++
		}
	}
	return i == len(r)
}

func main() {
	fmt.Println("Topic: longest common subsequence and edit distance")
	failed := false

	for _, tc := range []struct {
		a, b    string
		wantLen int
	}{
		{"ABCBDAB", "BDCABA", 4},
		{"abcde", "ace", 3},
		{"kitten", "sitting", 4},
		{"abc", "xyz", 0},
		{"नमस्ते", "नमन", 2},
	} {
		got := LCS(tc.a, tc.b)
		ok := len([]rune(got)) == tc.wantLen && isSubsequence(got, tc.a) && isSubsequence(got, tc.b)
		status := "ok"
		if !ok {
			status = "MISMATCH"
			failed = true
		}
		fmt.Printf("  LCS(%q, %q) = %q (len %d) %s\n", tc.a, tc.b, got, len([]rune(got)), status)
	}

	for _, tc := range []struct {
		a, b string
		want int
	}{{"kitten", "sitting", 3}, {"flaw", "lawn", 2}, {"", "abc", 3}, {"same", "same", 0}} {
		got := EditDistance(tc.a, tc.b)
		status := "ok"
		if got != tc.want {
			status = "MISMATCH"
			failed = true
		}
		fmt.Printf("  EditDistance(%q, %q) = %d %s\n", tc.a, tc.b, got, status)
	}

	if failed {
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}
