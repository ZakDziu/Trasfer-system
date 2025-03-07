package config

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// LoadTestConfig loads configuration for tests
func LoadTestConfig(t *testing.T) *Config {
	t.Helper()

	// Save current env vars
	oldEnv := map[string]string{
		"GO_ENV":      os.Getenv("GO_ENV"),
		"CONFIG_PATH": os.Getenv("CONFIG_PATH"),
		"DB_HOST":     os.Getenv("DB_HOST"),
		"DB_PORT":     os.Getenv("DB_PORT"),
		"DB_NAME":     os.Getenv("DB_NAME"),
		"DB_USER":     os.Getenv("DB_USER"),
		"DB_PASSWORD": os.Getenv("DB_PASSWORD"),
	}

	// Restore env vars after test
	t.Cleanup(func() {
		for k, v := range oldEnv {
			os.Setenv(k, v)
		}
	})

	// Set test env vars
	os.Setenv("GO_ENV", "test")
	if os.Getenv("CONFIG_PATH") == "" {
		// Try to find project root
		currentDir, err := os.Getwd()
		require.NoError(t, err)

		for dir := currentDir; dir != "/"; dir = filepath.Dir(dir) {
			if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
				os.Setenv("CONFIG_PATH", dir)
				break
			}
		}
	}

	// Create test database if it doesn't exist
	cfg, err := Load()
	require.NoError(t, err)

	// Connect to default database to create test database
	defaultDB := cfg.Database
	defaultDB.DBName = "postgres" // connect to default postgres database
	db, err := sql.Open("postgres", defaultDB.GetDSN())
	require.NoError(t, err)
	defer db.Close()

	// Create test database if it doesn't exist
	_, err = db.Exec("CREATE DATABASE money_transfer_test")
	if err != nil {
		// Ignore error if database already exists
		if !strings.Contains(err.Error(), "already exists") {
			require.NoError(t, err)
		}
	}

	return cfg
}
