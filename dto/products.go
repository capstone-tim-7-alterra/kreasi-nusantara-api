package dto

import (
	"time"

	"github.com/google/uuid"
)

type ProductResponse struct {
	ID              uuid.UUID `json:"id"`
	Image           string    `json:"image"`
	Name            string    `json:"name"`
	OriginalPrice   int       `json:"original_price"`
	DiscountPercent *int      `json:"discount_percent"`
	DiscountPrice   *float64  `json:"discount_price"`
	AverageRating   float64   `json:"average_rating"`
	TotalReview     int       `json:"total_review"`
}

type ProductDetailResponse struct {
	ID              uuid.UUID                `json:"id"`
	Name            string                   `json:"name"`
	Description     string                   `json:"description"`
	Images          []string                 `json:"images"`
	Videos          []string                 `json:"videos"`
	OriginalPrice   int                      `json:"original_price"`
	DiscountPercent *int                     `json:"discount_percent,omitempty"`
	DiscountPrice   *float64                 `json:"discount_price,omitempty"`
	AverageRating   float64                  `json:"average_rating"`
	TotalReview     int                      `json:"total_review"`
	LatestReview    []*ProductReviewResponse   `json:"latest_review,omitempty"`
	Variants        []ProductVariantResponse `json:"variants"`
}

type ProductReviewRequest struct {
	Rating int    `json:"rating"`
	Review string `json:"review"`
}

type ProductReviewResponse struct {
	User      UserReview `json:"user"`
	Rating    int        `json:"rating"`
	CreatedAt time.Time  `json:"created_at"`
	Review    string     `json:"review"`
}

type ProductVariantResponse struct {
	Size  string `json:"size"`
	Stock int    `json:"stock"`
}

type UserReview struct {
	ImageURL *string `json:"image_url"`
	Username string  `json:"username"`
}
