package postgres

import (
	"context"
	"database/sql"
	"log"

	"money-transfer/internal/domain/models"
	"money-transfer/internal/domain/transfer_errors"
)

// AccountRepository handles all database operations related to accounts
type AccountRepository struct {
	db *sql.DB
}

// NewAccountRepository creates a new instance of AccountRepository
func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

// GetAccount retrieves account information by ID
func (r *AccountRepository) GetAccount(ctx context.Context, id string) (*models.Account, error) {
	var account models.Account
	err := r.db.QueryRowContext(ctx, "SELECT id, balance FROM accounts WHERE id = $1", id).
		Scan(&account.ID, &account.Balance)

	if err == sql.ErrNoRows {
		return nil, transfererrors.ErrAccountNotFound
	}
	if err != nil {
		return nil, err
	}

	return &account, nil
}

// InitializeTestData populates the database with test accounts
func (r *AccountRepository) InitializeTestData(ctx context.Context) error {
	accounts := []struct {
		id      string
		balance float64
	}{
		{"Mark", 100},
		{"Jane", 50},
		{"Adam", 0},
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Printf("Error rolling back transaction: %v", err)
		}
	}()

	for _, acc := range accounts {
		_, err := tx.ExecContext(ctx,
			"INSERT INTO accounts (id, balance) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET balance = $2",
			acc.id, acc.balance)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// TransferWithinTx performs a money transfer between accounts within a transaction
// Uses serializable isolation level to prevent concurrent modifications
func (r *AccountRepository) TransferWithinTx(ctx context.Context, fromID, toID string, amount float64) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	defer func() {
		if err = tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Printf("Error rolling back transaction: %v", err)
		}
	}()

	result, err := tx.ExecContext(ctx, `
		UPDATE accounts 
		SET balance = balance - $1 
		WHERE id = $2 AND balance >= $1`,
		amount, fromID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		var exists bool
		err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM accounts WHERE id = $1)", fromID).Scan(&exists)
		if err != nil {
			return err
		}
		if !exists {
			return transfererrors.ErrAccountNotFound
		}
		return transfererrors.ErrInsufficientFunds
	}

	result, err = tx.ExecContext(ctx,
		"UPDATE accounts SET balance = balance + $1 WHERE id = $2",
		amount, toID)
	if err != nil {
		return err
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return transfererrors.ErrAccountNotFound
	}

	return tx.Commit()
}
