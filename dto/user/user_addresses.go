package dto

import "github.com/google/uuid"

type UserAddressRequest struct {
	Label         string `json:"label" validate:"required"`
	RecipientName string `json:"recipient_name" validate:"required"`
	Phone         string `json:"phone" validate:"required"`
	Address       string `json:"address" validate:"required"`
	City          string `json:"city" validate:"required"`
	Province      string `json:"province" validate:"required"`
	PostalCode    string `json:"postal_code" validate:"required"`
	IsPrimary     bool   `json:"is_primary"`
}

type UserAddressResponse struct {
	ID            uuid.UUID `json:"id"`
	Label         string    `json:"label" validate:"required"`
	RecipientName string    `json:"recipient_name" validate:"required"`
	Phone         string    `json:"phone" validate:"required"`
	Address       string    `json:"address" validate:"required"`
	City          string    `json:"city" validate:"required"`
	Province      string    `json:"province" validate:"required"`
	PostalCode    string    `json:"postal_code" validate:"required"`
	IsPrimary     bool      `json:"is_primary"`
}
