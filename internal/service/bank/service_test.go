package bank

import (
	"context"
	"testing"

	"money-transfer/internal/domain/models"
	"money-transfer/internal/domain/transfer_errors"
	"money-transfer/internal/storage/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBankService_Transfer(t *testing.T) {
	tests := []struct {
		name    string
		req     models.TransferRequest
		mock    func(*mocks.Store, *mocks.AccountRepository)
		wantErr error
	}{
		{
			name: "successful transfer",
			req: models.TransferRequest{
				From:   "Mark",
				To:     "Jane",
				Amount: 50,
			},
			mock: func(s *mocks.Store, ar *mocks.AccountRepository) {
				s.On("Account").Return(ar)
				ar.On("TransferWithinTx",
					mock.Anything,
					"Mark",
					"Jane",
					float64(50),
				).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "same account transfer",
			req: models.TransferRequest{
				From:   "Mark",
				To:     "Mark",
				Amount: 50,
			},
			mock:    func(_ *mocks.Store, _ *mocks.AccountRepository) {},
			wantErr: transfererrors.ErrSameAccount,
		},
		{
			name: "negative amount",
			req: models.TransferRequest{
				From:   "Mark",
				To:     "Jane",
				Amount: -50,
			},
			mock:    func(_ *mocks.Store, _ *mocks.AccountRepository) {},
			wantErr: transfererrors.ErrInvalidAmount,
		},
		{
			name: "insufficient funds",
			req: models.TransferRequest{
				From:   "Mark",
				To:     "Jane",
				Amount: 50,
			},
			mock: func(s *mocks.Store, ar *mocks.AccountRepository) {
				s.On("Account").Return(ar)
				ar.On("TransferWithinTx",
					mock.Anything,
					"Mark",
					"Jane",
					float64(50),
				).Return(transfererrors.ErrInsufficientFunds)
			},
			wantErr: transfererrors.ErrInsufficientFunds,
		},
		{
			name: "account not found",
			req: models.TransferRequest{
				From:   "NonExistent",
				To:     "Jane",
				Amount: 50,
			},
			mock: func(s *mocks.Store, ar *mocks.AccountRepository) {
				s.On("Account").Return(ar)
				ar.On("TransferWithinTx",
					mock.Anything,
					"NonExistent",
					"Jane",
					float64(50),
				).Return(transfererrors.ErrAccountNotFound)
			},
			wantErr: transfererrors.ErrAccountNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			mockStore := mocks.NewStore(t)
			mockAccountRepo := mocks.NewAccountRepository(t)
			tt.mock(mockStore, mockAccountRepo)

			// Create service with mock
			service := NewService(mockStore)

			// Execute test
			err := service.Transfer(context.Background(), tt.req)

			// Check results
			assert.ErrorIs(t, err, tt.wantErr)
			mockStore.AssertExpectations(t)
			mockAccountRepo.AssertExpectations(t)
		})
	}
}

func TestBankService_GetBalance(t *testing.T) {
	tests := []struct {
		name        string
		accountID   string
		mock        func(*mocks.Store, *mocks.AccountRepository)
		wantBalance float64
		wantErr     error
	}{
		{
			name:      "successful balance retrieval",
			accountID: "Mark",
			mock: func(s *mocks.Store, ar *mocks.AccountRepository) {
				s.On("Account").Return(ar)
				ar.On("GetAccount",
					mock.Anything,
					"Mark",
				).Return(&models.Account{
					ID:      "Mark",
					Balance: 100,
				}, nil)
			},
			wantBalance: 100,
			wantErr:     nil,
		},
		{
			name:      "account not found",
			accountID: "NonExistent",
			mock: func(s *mocks.Store, ar *mocks.AccountRepository) {
				s.On("Account").Return(ar)
				ar.On("GetAccount",
					mock.Anything,
					"NonExistent",
				).Return(nil, transfererrors.ErrAccountNotFound)
			},
			wantBalance: 0,
			wantErr:     transfererrors.ErrAccountNotFound,
		},
		{
			name:      "database error",
			accountID: "Mark",
			mock: func(s *mocks.Store, ar *mocks.AccountRepository) {
				s.On("Account").Return(ar)
				ar.On("GetAccount",
					mock.Anything,
					"Mark",
				).Return(nil, assert.AnError)
			},
			wantBalance: 0,
			wantErr:     assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			mockStore := mocks.NewStore(t)
			mockAccountRepo := mocks.NewAccountRepository(t)
			tt.mock(mockStore, mockAccountRepo)

			// Create service with mock
			service := NewService(mockStore)

			// Execute test
			balance, err := service.GetBalance(context.Background(), tt.accountID)

			// Check results
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantBalance, balance)
			mockStore.AssertExpectations(t)
			mockAccountRepo.AssertExpectations(t)
		})
	}
}
