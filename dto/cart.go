package dto

import (
	"time"

	"github.com/google/uuid"
)

type CartItemResponse struct {
	ID        uuid.UUID            `json:"id"`
	Products  []ProductInformation `json:"products"`
	Size      string               `json:"size"`
	Quantity  int                  `json:"quantity"`
	CreatedAt time.Time            `json:"created_at"`
}

type AddCartItemRequest struct {
	ProductVariantID uuid.UUID `json:"product_variant_id" form:"product_variant_id"`
	Quantity         int       `json:"quantity" form:"quantity"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" form:"quantity"`
}

type ProductInformation struct {
	ProductID     uuid.UUID `json:"product_id"`
	ProductName   string    `json:"product_name"`
	ProductImage  string    `json:"product_image"`
	OriginalPrice int       `json:"original_price"`
	DiscountPrice float64   `json:"discount_price,omitempty"`
}
