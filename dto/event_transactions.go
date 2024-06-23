package dto

import "github.com/google/uuid"

type EventTransactionRequest struct {
	EventPriceID   uuid.UUID `json:"event_price_id" validate:"required"`
	Quantity       int       `json:"quantity" validate:"required"`
	IdentityNumber string    `json:"identity_number" validate:"required"`
	FullName       string    `json:"full_name" validate:"required"`
	Email          string    `json:"email" validate:"required,email"`
	Phone          string    `json:"phone" validate:"required"`
}

type EventTransactionResponse struct {
	ID                uuid.UUID        `json:"id"`
	EventPriceID      uuid.UUID        `json:"event_price_id"`
	UserID            uuid.UUID        `json:"user_id"`
	BuyerInformation  BuyerInformation `json:"buyer_information"`
	TotalAmount       float64          `json:"total_amount"`
	TransactionStatus string           `json:"transaction_status"`
	SnapURL           string           `json:"snap_url"`
}

type BuyerInformation struct {
	IdentityNumber string `json:"identity_number" validate:"required"`
	FullName       string `json:"full_name" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	Phone          string `json:"phone" validate:"required"`
}
