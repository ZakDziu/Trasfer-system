package postgres

import (
	"context"
	"math"
	"sync"
	"testing"
	"time"

	"money-transfer/config"
	"money-transfer/internal/domain/transfer_errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *AccountRepository {
	t.Helper()
	cfg := config.LoadTestConfig(t)

	store, err := NewStore(cfg.Database.GetDSN())
	require.NoError(t, err)

	_, err = store.db.Exec(`
		SELECT pg_terminate_backend(pid) 
		FROM pg_stat_activity 
		WHERE datname = current_database() AND pid <> pg_backend_pid()
	`)
	require.NoError(t, err)

	_, err = store.db.Exec("TRUNCATE TABLE accounts")
	require.NoError(t, err)

	return store.accountRepo.(*AccountRepository)
}

func TestAccountRepository_GetAccount(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	// Initialize test data
	err := repo.InitializeTestData(ctx)
	require.NoError(t, err)

	tests := []struct {
		name          string
		accountID     string
		expectedError error
		checkBalance  float64
	}{
		{
			name:         "successful account retrieval",
			accountID:    "Mark",
			checkBalance: 100,
		},
		{
			name:          "account not found",
			accountID:     "NonExistent",
			expectedError: transfererrors.ErrAccountNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account, err := repo.GetAccount(ctx, tt.accountID)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, account)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, account)
				assert.Equal(t, tt.accountID, account.ID)
				assert.Equal(t, tt.checkBalance, account.Balance)
			}
		})
	}
}

func TestAccountRepository_TransferWithinTx(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	err := repo.InitializeTestData(ctx)
	require.NoError(t, err)

	tests := []struct {
		name          string
		fromID        string
		toID          string
		amount        float64
		expectedError error
	}{
		{
			name:   "successful transfer",
			fromID: "Mark",
			toID:   "Jane",
			amount: 50,
		},
		{
			name:          "insufficient funds",
			fromID:        "Adam",
			toID:          "Jane",
			amount:        50,
			expectedError: transfererrors.ErrInsufficientFunds,
		},
		{
			name:          "sender does not exist",
			fromID:        "NonExistent",
			toID:          "Jane",
			amount:        50,
			expectedError: transfererrors.ErrAccountNotFound,
		},
		{
			name:          "recipient does not exist",
			fromID:        "Mark",
			toID:          "NonExistent",
			amount:        50,
			expectedError: transfererrors.ErrAccountNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get initial balances
			var fromBalanceBefore, toBalanceBefore float64
			fromAcc, getErr := repo.GetAccount(ctx, tt.fromID)
			if getErr == nil {
				fromBalanceBefore = fromAcc.Balance
			}
			toAcc, getErr := repo.GetAccount(ctx, tt.toID)
			if getErr == nil {
				toBalanceBefore = toAcc.Balance
			}

			// Perform transfer
			err := repo.TransferWithinTx(ctx, tt.fromID, tt.toID, tt.amount)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)

				// Check that balances haven't changed
				fromAcc, getErr := repo.GetAccount(ctx, tt.fromID)
				if getErr == nil {
					assert.Equal(t, fromBalanceBefore, fromAcc.Balance)
				}
				toAcc, getErr := repo.GetAccount(ctx, tt.toID)
				if getErr == nil {
					assert.Equal(t, toBalanceBefore, toAcc.Balance)
				}
			} else {
				assert.NoError(t, err)

				// Check new balances
				fromAcc, getErr := repo.GetAccount(ctx, tt.fromID)
				require.NoError(t, getErr)
				assert.Equal(t, fromBalanceBefore-tt.amount, fromAcc.Balance)

				toAcc, getErr := repo.GetAccount(ctx, tt.toID)
				require.NoError(t, getErr)
				assert.Equal(t, toBalanceBefore+tt.amount, toAcc.Balance)
			}
		})
	}
}

func TestAccountRepository_ConcurrentTransfers(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	err := repo.InitializeTestData(ctx)
	require.NoError(t, err)

	numTransfers := 10
	transferAmount := 1.0

	markAcc, err := repo.GetAccount(ctx, "Mark")
	require.NoError(t, err)
	janeAcc, err := repo.GetAccount(ctx, "Jane")
	require.NoError(t, err)

	initialMarkBalance := markAcc.Balance
	initialJaneBalance := janeAcc.Balance

	var wg sync.WaitGroup
	var mu sync.Mutex
	var successfulTransfers int

	for i := 0; i < numTransfers; i++ {
		wg.Add(2)

		// Transfer from Mark to Jane
		go func() {
			defer wg.Done()
			// Use a separate error variable to avoid shadowing
			if transferErr := repo.TransferWithinTx(ctx, "Mark", "Jane", transferAmount); transferErr == nil {
				mu.Lock()
				successfulTransfers++
				mu.Unlock()
			}
		}()

		// Transfer from Jane to Mark
		go func() {
			defer wg.Done()
			// Use a separate error variable to avoid shadowing
			if transferErr := repo.TransferWithinTx(ctx, "Jane", "Mark", transferAmount); transferErr == nil {
				mu.Lock()
				successfulTransfers++
				mu.Unlock()
			}
		}()

		time.Sleep(10 * time.Millisecond)
	}

	wg.Wait()

	// Check final balances
	markAcc, err = repo.GetAccount(ctx, "Mark")
	require.NoError(t, err)
	janeAcc, err = repo.GetAccount(ctx, "Jane")
	require.NoError(t, err)

	t.Logf("Successful transfers: %d out of %d attempts", successfulTransfers, numTransfers*2)
	t.Logf("Final balances - Mark: %.2f, Jane: %.2f", markAcc.Balance, janeAcc.Balance)

	totalBalanceBefore := initialMarkBalance + initialJaneBalance
	totalBalanceAfter := markAcc.Balance + janeAcc.Balance
	assert.Equal(t, totalBalanceBefore, totalBalanceAfter,
		"Total balance changed: before=%.2f, after=%.2f", totalBalanceBefore, totalBalanceAfter)

	maxBalanceChange := float64(successfulTransfers) * transferAmount
	assert.LessOrEqual(t, math.Abs(markAcc.Balance-initialMarkBalance), maxBalanceChange,
		"Mark's balance changed more than expected: initial=%.2f, final=%.2f, max change=%.2f",
		initialMarkBalance, markAcc.Balance, maxBalanceChange)
	assert.LessOrEqual(t, math.Abs(janeAcc.Balance-initialJaneBalance), maxBalanceChange,
		"Jane's balance changed more than expected: initial=%.2f, final=%.2f, max change=%.2f",
		initialJaneBalance, janeAcc.Balance, maxBalanceChange)
}
