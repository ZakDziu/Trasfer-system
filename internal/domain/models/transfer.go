package models

// TransferRequest represents the input data for a money transfer operation
type TransferRequest struct {
	From   string  `json:"from"`   // Source account ID
	To     string  `json:"to"`     // Destination account ID
	Amount float64 `json:"amount"` // Amount to transfer
}

// TransferResponse represents the result of a transfer operation
type TransferResponse struct {
	Success bool   `json:"success"`           // Indicates if transfer was successful
	Message string `json:"message,omitempty"` // Optional error or success message
}
