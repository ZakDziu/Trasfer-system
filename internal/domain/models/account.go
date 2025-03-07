package models

// Account represents a bank account entity
// ID is a unique identifier for the account
// Balance represents the current monetary amount in the account
type Account struct {
	ID      string  `json:"id"`
	Balance float64 `json:"balance"`
}
