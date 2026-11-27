// Topic: n-queens | Concept: backtracking | Stage: 04-dsa
// Run:   go run ./04-dsa/backtracking/n-queens
//
// Explore:
// * Place n queens on an n x n board so no two attack each other
// * One queen per row: for each row try every column, skip attacked squares, recurse
// * O(1) attack checks with three sets: columns, "/" diagonals (r+c), "\" diagonals (r-c)
// * Pruning (skipping attacked squares early) is what makes backtracking fast
// * Worst case is still exponential, roughly O(n!)
// * Same pattern: Sudoku solvers, maze solving, scheduling with constraints
package main

import (
	"fmt"
	"os"
	"strings"
)

// SolveNQueens returns every solution as a slice of column indexes, one per row.
func SolveNQueens(n int) [][]int {
	var solutions [][]int
	cols := make([]bool, n)
	diag1 := make([]bool, 2*n) // r + c
	diag2 := make([]bool, 2*n) // r - c + n
	queens := make([]int, n)

	var place func(row int)
	place = func(row int) {
		if row == n {
			solutions = append(solutions, append([]int(nil), queens...))
			return
		}
		for c := range n {
			if cols[c] || diag1[row+c] || diag2[row-c+n] {
				continue // attacked: prune this branch
			}
			queens[row] = c
			cols[c], diag1[row+c], diag2[row-c+n] = true, true, true
			place(row + 1)
			cols[c], diag1[row+c], diag2[row-c+n] = false, false, false
		}
	}
	place(0)
	return solutions
}

func board(sol []int) string {
	var b strings.Builder
	for _, c := range sol {
		row := []byte(strings.Repeat(". ", len(sol)))
		row[2*c] = 'Q'
		b.WriteString("    " + string(row) + "\n")
	}
	return b.String()
}

// safe double-checks a finished board the slow way.
func safe(sol []int) bool {
	for r1 := range sol {
		for r2 := r1 + 1; r2 < len(sol); r2++ {
			dc := sol[r1] - sol[r2]
			if dc == 0 || dc == r2-r1 || dc == r1-r2 {
				return false
			}
		}
	}
	return true
}

func main() {
	fmt.Println("Topic: N-Queens")
	failed := false

	known := map[int]int{1: 1, 2: 0, 3: 0, 4: 2, 5: 10, 6: 4, 7: 40, 8: 92}
	for n := 1; n <= 8; n++ {
		sols := SolveNQueens(n)
		status := "ok"
		if len(sols) != known[n] {
			status = fmt.Sprintf("MISMATCH (want %d)", known[n])
			failed = true
		}
		for _, s := range sols {
			if !safe(s) {
				status = "UNSAFE BOARD"
				failed = true
			}
		}
		fmt.Printf("  n=%d: %2d solutions %s\n", n, len(sols), status)
	}

	first := SolveNQueens(8)[0]
	fmt.Println("first 8-queens solution:", first)
	fmt.Print(board(first))

	if failed {
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}
