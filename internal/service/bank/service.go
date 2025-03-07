package bank

import (
	"context"
	"log"

	"money-transfer/internal/domain/models"
	"money-transfer/internal/domain/transfer_errors"
	"money-transfer/internal/storage"
)

// Service handles all banking operations
type Service struct {
	store storage.Store
}

// NewService creates a new instance of banking service
func NewService(store storage.Store) *Service {
	return &Service{
		store: store,
	}
}

// Transfer performs a money transfer between two accounts
// Returns error if transfer cannot be completed
func (s *Service) Transfer(ctx context.Context, req models.TransferRequest) error {
	log.Printf("Transfer request: %+v", req)
	if req.From == req.To {
		return transfererrors.ErrSameAccount
	}

	if req.Amount <= 0 {
		return transfererrors.ErrInvalidAmount
	}

	err := s.store.Account().TransferWithinTx(ctx, req.From, req.To, req.Amount)
	if err != nil {
		log.Printf("Transfer failed: %v", err)
		return err
	}

	return nil
}

// GetBalance returns the current balance for the specified account
// Returns error if account cannot be found
func (s *Service) GetBalance(ctx context.Context, accountID string) (float64, error) {
	account, err := s.store.Account().GetAccount(ctx, accountID)
	if err != nil {
		return 0, err
	}

	return account.Balance, nil
}
