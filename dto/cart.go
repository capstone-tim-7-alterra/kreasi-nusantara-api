package dto

import (
	"time"

	"github.com/google/uuid"
)

type CartItemResponse struct {
	ID        uuid.UUID            `json:"id"`
	Products  []ProductInformation `json:"products"`
	Total     float64              `json:"total"`
	CreatedAt time.Time            `json:"created_at"`
}

type AddCartItemRequest struct {
	ProductVariantID uuid.UUID `json:"product_variant_id" form:"product_variant_id"`
	Quantity         int       `json:"quantity" form:"quantity"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" form:"quantity"`
}

type RemoveCartItemRequest struct {
	ProductVariantID uuid.UUID `json:"product_variant_id" form:"product_variant_id"`
}

type ProductInformation struct {
	CartItemID       uuid.UUID `json:"cart_item_id" form:"cart_item_id"`
	CartID           uuid.UUID `json:"cart_id"`
	ProductVariantID uuid.UUID `json:"product_variant_id"`
	ProductName      string    `json:"product_name"`
	ProductImage     string    `json:"product_image"`
	OriginalPrice    int       `json:"original_price"`
	DiscountPrice    float64   `json:"discount_price,omitempty"`
	Size             string    `json:"size"`
	Quantity         int       `json:"quantity"`
}
