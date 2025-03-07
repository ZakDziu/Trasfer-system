package mocks

import (
	"context"
	"money-transfer/internal/domain/models"

	"github.com/stretchr/testify/mock"
)

type BankServiceMock struct {
	mock.Mock
}

func (m *BankServiceMock) Transfer(ctx context.Context, req models.TransferRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *BankServiceMock) GetBalance(ctx context.Context, accountID string) (float64, error) {
	args := m.Called(ctx, accountID)
	return args.Get(0).(float64), args.Error(1)
}
