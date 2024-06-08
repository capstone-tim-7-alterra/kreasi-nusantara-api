package dto

import "github.com/google/uuid"

type ProductResponse struct {
	ID              uuid.UUID `json:"id"`
	Image           string    `json:"image"`
	ProductName     string    `json:"product_name"`
	Price           int   `json:"product_price"`
	Rating          *float64  `json:"rating"`
	NumberOfReviews *int      `json:"number_of_reviews"`
}

type ProductDetailResponse struct {
	ID          uuid.UUID `json:"id"`
	ProductName string    `json:"product_name"`
	Description string    `json:"description"`
}
