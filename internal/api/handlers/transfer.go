package handlers

import (
	"errors"
	"net/http"

	"money-transfer/internal/domain/models"
	"money-transfer/internal/domain/transfer_errors"
	"money-transfer/internal/service"

	"github.com/gin-gonic/gin"
)

// TransferHandler handles money transfer requests
type TransferHandler struct {
	bankService service.BankService
}

// NewTransferHandler creates a new transfer handler
func NewTransferHandler(cfg *HandlerConfig) *TransferHandler {
	return &TransferHandler{
		bankService: cfg.BankService,
	}
}

// Register registers handler routes
func (h *TransferHandler) Register(group *gin.RouterGroup) {
	group.POST("/transfer", h.Transfer)
}

// Transfer godoc
// @Summary Execute money transfer between accounts
// @Description Transfers specified amount from one account to another
// @Tags transfer
// @Accept json
// @Produce json
// @Param request body models.TransferRequest true "Transfer details"
// @Success 200 {object} models.TransferResponse "Successful transfer"
// @Failure 400 {object} map[string]string "Validation error"
// @Failure 404 {object} map[string]string "Account not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /transfer [post]
func (h *TransferHandler) Transfer(c *gin.Context) {
	var req models.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.bankService.Transfer(c.Request.Context(), req); err != nil {
		switch {
		case errors.Is(err, transfererrors.ErrAccountNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, transfererrors.ErrInsufficientFunds),
			errors.Is(err, transfererrors.ErrInvalidAmount),
			errors.Is(err, transfererrors.ErrSameAccount):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, models.TransferResponse{Success: true})
}
