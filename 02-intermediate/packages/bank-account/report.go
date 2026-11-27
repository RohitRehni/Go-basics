package main

// A second file in package main. Everything here is visible to main.go
// without an import, because both files are in the same folder (package).

import (
	"fmt"
	"strings"

	"learning/02-intermediate/packages/bank-account/account"
)

func printStatement(a *account.Account) {
	line := strings.Repeat("-", 40)
	fmt.Println(line)
	fmt.Printf("Statement for %s\n", a.Owner)
	fmt.Println(line)
	for _, entry := range a.History() {
		fmt.Println(entry)
	}
	fmt.Println(line)
}
