package postgres

import (
	"database/sql"

	"money-transfer/internal/storage"

	// Import PostgreSQL driver for side effects - registers postgres driver
	_ "github.com/lib/pq"
)

// Store implements the Store interface for PostgreSQL database
type Store struct {
	db          *sql.DB
	accountRepo storage.AccountRepository
}

// NewStore creates a new instance of Store and initializes the database
// Returns error if database connection or schema creation fails
func NewStore(connStr string) (*Store, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := createSchema(db); err != nil {
		return nil, err
	}

	store := &Store{
		db: db,
	}
	store.accountRepo = NewAccountRepository(db)

	return store, nil
}

// createSchema ensures that the required database tables exist
func createSchema(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS accounts (
			id VARCHAR(255) PRIMARY KEY,
			balance DECIMAL(10, 2) NOT NULL
		)`

	_, err := db.Exec(query)
	return err
}

// DB returns the underlying database connection
func (s *Store) DB() *sql.DB {
	return s.db
}

// Account returns the account repository instance
func (s *Store) Account() storage.AccountRepository {
	return s.accountRepo
}
