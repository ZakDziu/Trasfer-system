// Package transfererrors provides error definitions for transfer operations
package transfererrors

import "errors"

// Common errors that can occur during transfer operations
var (
	// ErrAccountNotFound is returned when the specified account doesn't exist
	ErrAccountNotFound = errors.New("account not found")

	// ErrInsufficientFunds is returned when the source account has insufficient balance
	ErrInsufficientFunds = errors.New("insufficient funds")

	// ErrInvalidAmount is returned when the transfer amount is invalid (e.g., negative or zero)
	ErrInvalidAmount = errors.New("invalid amount")

	// ErrSameAccount is returned when trying to transfer money to the same account
	ErrSameAccount = errors.New("cannot transfer to same account")
)
