package handlers

import "money-transfer/internal/service"

// Factory creates and manages all handlers
type Factory struct {
	config *HandlerConfig
}

// NewFactory creates a new handler factory
func NewFactory(bankService service.BankService) *Factory {
	return &Factory{
		config: NewHandlerConfig(bankService),
	}
}

// CreateHandlers creates all application handlers
func (f *Factory) CreateHandlers() []Handler {
	return []Handler{
		NewTransferHandler(f.config),
		NewBalanceHandler(f.config),
	}
}
