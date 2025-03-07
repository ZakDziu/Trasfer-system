package bank

import (
	"context"
	"fmt"
	"os"
	"testing"

	"money-transfer/config"
	"money-transfer/internal/domain/models"
	"money-transfer/internal/storage/postgres"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testStore *postgres.Store

func TestMain(m *testing.M) {
	// Load test configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Connect to test database
	testStore, err = postgres.NewStore(cfg.Database.GetDSN())
	if err != nil {
		fmt.Printf("Failed to connect to test database: %v\n", err)
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Cleanup test DB
	if err := cleanupTestDB(); err != nil {
		fmt.Printf("Failed to cleanup test database: %v\n", err)
	}

	os.Exit(code)
}

func cleanupTestDB() error {
	_, err := testStore.DB().Exec(`
		SELECT pg_terminate_backend(pid) 
		FROM pg_stat_activity 
		WHERE datname = current_database() AND pid <> pg_backend_pid()
	`)
	if err != nil {
		return err
	}

	_, err = testStore.DB().Exec("TRUNCATE accounts")
	return err
}

func setupTest(t *testing.T) *Service {
	t.Helper()

	// Clean table before each test
	require.NoError(t, cleanupTestDB())

	// Initialize test data
	ctx := context.Background()
	require.NoError(t, testStore.Account().InitializeTestData(ctx))

	return NewService(testStore)
}

func TestBankService_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	service := setupTest(t)
	ctx := context.Background()

	// Test successful transfer
	err := service.Transfer(ctx, models.TransferRequest{
		From:   "Mark",
		To:     "Jane",
		Amount: 50,
	})
	require.NoError(t, err)

	// Verify balances after transfer
	markBalance, err := service.GetBalance(ctx, "Mark")
	require.NoError(t, err)
	assert.Equal(t, 50.0, markBalance)

	janeBalance, err := service.GetBalance(ctx, "Jane")
	require.NoError(t, err)
	assert.Equal(t, 100.0, janeBalance)
}
