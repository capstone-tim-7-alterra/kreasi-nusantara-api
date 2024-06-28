package dto

import "github.com/google/uuid"

type TransactionRequest struct {
	CartId uuid.UUID `json:"cart_id" validate:"required"`
}

type SingleTransactionRequest struct {
	ProductVariantID uuid.UUID `json:"product_variant_id" validate:"required"`
	Quantity         int       `json:"quantity" validate:"required"`
}

type SingleTransactionResponse struct {
	ID                string      `json:"id"`
	UserID            uuid.UUID   `json:"user_id"`
	ProductInfo       ProductInfo `json:"product_info"`
	TotalAmount       float64     `json:"total_amount"`
	TransactionStatus string      `json:"transaction_status"`
	SnapURL           string      `json:"snap_url"`
}

type TransactionResponse struct {
	ID                string    `json:"id"`
	CartId            uuid.UUID `json:"cart_id"`
	UserId            uuid.UUID `json:"user_id"`
	TotalAmount       float64   `json:"total_amount"`
	TransactionStatus string    `json:"transaction_status"`
	SnapURL           string    `json:"snap_url"`
}

type ProductInfo struct {
	ProductVariantID uuid.UUID `json:"product_variant_id"`
	Quantity         int       `json:"quantity"`
}
