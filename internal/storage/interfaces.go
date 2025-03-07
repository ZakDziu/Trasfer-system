package storage

import (
	"context"
	"database/sql"

	"money-transfer/internal/domain/models"
)

// Store represents the main interface for database operations
type Store interface {
	DB() *sql.DB
	Account() AccountRepository
}

// AccountRepository defines the interface for account-related database operations
type AccountRepository interface {
	// GetAccount retrieves account information by ID
	GetAccount(ctx context.Context, id string) (*models.Account, error)

	// TransferWithinTx performs a money transfer between accounts
	TransferWithinTx(ctx context.Context, fromID, toID string, amount float64) error

	// InitializeTestData sets up test data in the database
	InitializeTestData(ctx context.Context) error
}
