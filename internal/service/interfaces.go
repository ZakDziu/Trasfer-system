package service

import (
	"context"

	"money-transfer/internal/domain/models"
)

type BankService interface {
	Transfer(ctx context.Context, req models.TransferRequest) error
	GetBalance(ctx context.Context, accountID string) (float64, error)
}
