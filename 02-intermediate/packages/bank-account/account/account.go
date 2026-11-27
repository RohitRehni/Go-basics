// Package account is the "class" of this lesson: a type, a constructor and
// methods, with the internal state hidden from other packages.
//
// Visibility in Go is decided by the first letter:
//   - Capitalised names (Account, NewAccount, Deposit, Owner) are exported
//   - lowercase names (balance, history, record) stay inside this package
package account

import (
	"errors"
	"fmt"
)

// Exported sentinel errors so callers can check them with errors.Is
var (
	ErrInvalidAmount     = errors.New("amount must be positive")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

// Account holds an owner (exported) and a balance (unexported).
// Code outside this package can read Owner but must use methods for balance.
type Account struct {
	Owner   string
	balance int64 // in paise/cents, never float for money
	history []string
}

// NewAccount is the constructor: Go has no `new Account()` with logic,
// so a NewXxx function is the convention.
func NewAccount(owner string, opening int64) (*Account, error) {
	if opening < 0 {
		return nil, fmt.Errorf("opening balance %d: %w", opening, ErrInvalidAmount)
	}
	a := &Account{Owner: owner, balance: opening}
	a.record("open", opening)
	return a, nil
}

// Deposit adds money. Pointer receiver because it changes the account.
func (a *Account) Deposit(amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("deposit %d: %w", amount, ErrInvalidAmount)
	}
	a.balance += amount
	a.record("deposit", amount)
	return nil
}

// Withdraw removes money if there is enough.
func (a *Account) Withdraw(amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("withdraw %d: %w", amount, ErrInvalidAmount)
	}
	if amount > a.balance {
		return fmt.Errorf("withdraw %d with balance %d: %w", amount, a.balance, ErrInsufficientFunds)
	}
	a.balance -= amount
	a.record("withdraw", amount)
	return nil
}

// Balance is a getter: read-only access to the unexported field.
func (a *Account) Balance() int64 {
	return a.balance
}

// History returns a copy so callers cannot change the internal slice.
func (a *Account) History() []string {
	return append([]string(nil), a.history...)
}

// record is unexported: only this package can call it.
func (a *Account) record(kind string, amount int64) {
	a.history = append(a.history, fmt.Sprintf("%-8s %6d -> balance %6d", kind, amount, a.balance))
}
