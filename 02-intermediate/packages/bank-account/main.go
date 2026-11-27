// Topic: bank-account | Concept: packages | Stage: 02-intermediate
// Run:   go run ./02-intermediate/packages/bank-account
//
// Multi-file lesson, Go's version of object-oriented code:
//
//	main.go             package main: uses the account package
//	report.go           package main too: a second file in the same package
//	account/account.go  package account: the type, constructor and methods
//
// Explore:
// * A folder is a package; every .go file in it shares the package's names
// * Import path = module path + folder: learning/02-intermediate/packages/bank-account/account
// * Exported (Capital) vs unexported (lowercase) names across packages
// * Constructor functions (NewAccount) instead of classes
// * Methods with pointer receivers to change state
// * Sentinel errors exported by the package, checked with errors.Is
package main

import (
	"errors"
	"fmt"

	"learning/02-intermediate/packages/bank-account/account"
)

func main() {
	fmt.Println("Topic: packages and multi-file programs")

	acc, err := account.NewAccount("Asha", 1000)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	_ = acc.Deposit(500)

	if err := acc.Withdraw(5000); err != nil {
		fmt.Println("withdraw failed:", err)
		fmt.Println("is ErrInsufficientFunds:", errors.Is(err, account.ErrInsufficientFunds))
	}

	if err := acc.Deposit(-10); err != nil {
		fmt.Println("deposit failed:", err)
	}

	_ = acc.Withdraw(300)

	fmt.Println("owner (exported field):", acc.Owner)
	fmt.Println("balance (via getter):", acc.Balance())
	// fmt.Println(acc.balance) // compile error: acc.balance undefined (cannot refer to unexported field balance)
	// acc.record("x", 1)       // compile error: acc.record undefined (cannot refer to unexported method record)

	printStatement(acc) // defined in report.go, same package main

	if _, err := account.NewAccount("Ravi", -1); err != nil {
		fmt.Println("constructor validation:", err)
	}
}
