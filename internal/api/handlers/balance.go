package handlers

import (
	"net/http"

	"money-transfer/internal/domain/transfer_errors"
	"money-transfer/internal/service"

	"github.com/gin-gonic/gin"
)

// BalanceHandler handles balance-related requests
type BalanceHandler struct {
	bankService service.BankService
}

// NewBalanceHandler creates a new balance handler
func NewBalanceHandler(cfg *HandlerConfig) *BalanceHandler {
	return &BalanceHandler{
		bankService: cfg.BankService,
	}
}

// Register registers handler routes
func (h *BalanceHandler) Register(group *gin.RouterGroup) {
	group.GET("/balance/:account", h.GetBalance)
}

// GetBalance godoc
// @Summary Get account balance
// @Description Returns the current balance of the specified account
// @Tags balance
// @Accept json
// @Produce json
// @Param account path string true "Account ID"
// @Success 200 {object} map[string]float64 "Successful response with balance"
// @Failure 404 {object} map[string]string "Account not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /balance/{account} [get]
func (h *BalanceHandler) GetBalance(c *gin.Context) {
	accountID := c.Param("account")

	balance, err := h.bankService.GetBalance(c.Request.Context(), accountID)
	if err != nil {
		switch err {
		case transfererrors.ErrAccountNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}
