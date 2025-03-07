package handlers

import (
	"money-transfer/internal/api/interfaces"
	"money-transfer/internal/service"
)

// HandlerConfig contains configuration for all handlers
type HandlerConfig struct {
	BankService service.BankService
}

// Handler represents a common interface for all handlers
type Handler = interfaces.Handler

// NewHandlerConfig creates a new handler configuration
func NewHandlerConfig(bankService service.BankService) *HandlerConfig {
	return &HandlerConfig{
		BankService: bankService,
	}
}
