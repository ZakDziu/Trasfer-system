package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"money-transfer/internal/api/testutil"
	"money-transfer/internal/domain/models"
	"money-transfer/internal/domain/transfer_errors"
	"money-transfer/internal/service/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupRouter(bankService *mocks.BankServiceMock) *gin.Engine {
	handlersFactory := NewFactory(bankService)
	appHandlers := handlersFactory.CreateHandlers()
	return testutil.SetupTestRouter(appHandlers)
}

func TestTransferHandler_Transfer(t *testing.T) {
	tests := []struct {
		name       string
		request    models.TransferRequest
		setupMock  func(*mocks.BankServiceMock)
		wantStatus int
		wantError  string
	}{
		{
			name: "successful transfer",
			request: models.TransferRequest{
				From:   "Mark",
				To:     "Jane",
				Amount: 50,
			},
			setupMock: func(m *mocks.BankServiceMock) {
				m.On("Transfer", mock.Anything, models.TransferRequest{
					From: "Mark", To: "Jane", Amount: 50,
				}).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "insufficient funds",
			request: models.TransferRequest{
				From:   "Adam",
				To:     "Jane",
				Amount: 100,
			},
			setupMock: func(m *mocks.BankServiceMock) {
				m.On("Transfer", mock.Anything, models.TransferRequest{
					From: "Adam", To: "Jane", Amount: 100,
				}).Return(transfererrors.ErrInsufficientFunds)
			},
			wantStatus: http.StatusBadRequest,
			wantError:  transfererrors.ErrInsufficientFunds.Error(),
		},
		{
			name: "account not found",
			request: models.TransferRequest{
				From:   "NonExistent",
				To:     "Jane",
				Amount: 50,
			},
			setupMock: func(m *mocks.BankServiceMock) {
				m.On("Transfer", mock.Anything, models.TransferRequest{
					From: "NonExistent", To: "Jane", Amount: 50,
				}).Return(transfererrors.ErrAccountNotFound)
			},
			wantStatus: http.StatusNotFound,
			wantError:  transfererrors.ErrAccountNotFound.Error(),
		},
		{
			name: "same account",
			request: models.TransferRequest{
				From:   "Mark",
				To:     "Mark",
				Amount: 50,
			},
			setupMock: func(m *mocks.BankServiceMock) {
				m.On("Transfer", mock.Anything, models.TransferRequest{
					From: "Mark", To: "Mark", Amount: 50,
				}).Return(transfererrors.ErrSameAccount)
			},
			wantStatus: http.StatusBadRequest,
			wantError:  transfererrors.ErrSameAccount.Error(),
		},
		{
			name: "invalid amount",
			request: models.TransferRequest{
				From:   "Mark",
				To:     "Jane",
				Amount: -50,
			},
			setupMock: func(m *mocks.BankServiceMock) {
				m.On("Transfer", mock.Anything, models.TransferRequest{
					From: "Mark", To: "Jane", Amount: -50,
				}).Return(transfererrors.ErrInvalidAmount)
			},
			wantStatus: http.StatusBadRequest,
			wantError:  transfererrors.ErrInvalidAmount.Error(),
		},
		{
			name: "internal error",
			request: models.TransferRequest{
				From:   "Mark",
				To:     "Jane",
				Amount: 50,
			},
			setupMock: func(m *mocks.BankServiceMock) {
				m.On("Transfer", mock.Anything, models.TransferRequest{
					From: "Mark", To: "Jane", Amount: 50,
				}).Return(assert.AnError)
			},
			wantStatus: http.StatusInternalServerError,
			wantError:  "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock service for each test
			mockService := new(mocks.BankServiceMock)
			tt.setupMock(mockService)

			router := setupRouter(mockService)

			body, err := json.Marshal(tt.request)
			require.NoError(t, err)

			req := httptest.NewRequest("POST", "/api/v1/transfer", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var response map[string]interface{}
			err = json.NewDecoder(w.Body).Decode(&response)
			require.NoError(t, err)

			if tt.wantError != "" {
				assert.Equal(t, tt.wantError, response["error"])
			} else {
				assert.Equal(t, true, response["success"])
			}

			// Verify that all expected mock calls were made
			mockService.AssertExpectations(t)
		})
	}
}

func TestBalanceHandler_GetBalance(t *testing.T) {
	tests := []struct {
		name        string
		accountID   string
		setupMock   func(*mocks.BankServiceMock)
		wantStatus  int
		wantBalance float64
		wantError   string
	}{
		{
			name:      "get existing account",
			accountID: "Mark",
			setupMock: func(m *mocks.BankServiceMock) {
				m.On("GetBalance", mock.Anything, "Mark").Return(100.0, nil)
			},
			wantStatus:  http.StatusOK,
			wantBalance: 100.0,
		},
		{
			name:      "get non-existing account",
			accountID: "NonExistent",
			setupMock: func(m *mocks.BankServiceMock) {
				m.On("GetBalance", mock.Anything, "NonExistent").Return(0.0, transfererrors.ErrAccountNotFound)
			},
			wantStatus: http.StatusNotFound,
			wantError:  transfererrors.ErrAccountNotFound.Error(),
		},
		{
			name:      "internal error",
			accountID: "Mark",
			setupMock: func(m *mocks.BankServiceMock) {
				m.On("GetBalance", mock.Anything, "Mark").Return(0.0, assert.AnError)
			},
			wantStatus: http.StatusInternalServerError,
			wantError:  "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock service for each test
			mockService := new(mocks.BankServiceMock)
			tt.setupMock(mockService)

			router := setupRouter(mockService)

			req := httptest.NewRequest("GET", "/api/v1/balance/"+tt.accountID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			require.NoError(t, err)

			if tt.wantError != "" {
				assert.Equal(t, tt.wantError, response["error"])
			} else {
				assert.Equal(t, tt.wantBalance, response["balance"])
			}

			// Verify that all expected mock calls were made
			mockService.AssertExpectations(t)
		})
	}
}
