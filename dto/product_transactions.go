package dto

import "github.com/google/uuid"

type TransactionRequest struct {
	CartId uuid.UUID `json:"cart_id" validate:"required"`
}

type TransactionResponse struct {
	ID                string    `json:"id"`
	CartId            uuid.UUID `json:"cart_id"`
	UserId            uuid.UUID `json:"user_id"`
	TotalAmount       float64   `json:"total_amount"`
	TransactionStatus string    `json:"transaction_status"`
	SnapURL           string    `json:"snap_url"`
}
