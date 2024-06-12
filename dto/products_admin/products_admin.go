package productsadmin

import (
	"github.com/google/uuid"
)

type ProductRequest struct {
	Name            string                    `json:"name" form:"name" validate:"required"`
	Description     string                    `json:"description" form:"description" validate:"required"`
	MinOrder        int                       `json:"min_order" form:"min_order" validate:"required"`
	ProductPricing  ProductPricingRequest     `json:"product_pricing" form:"product_pricing" validate:"required"`
	CategoryID      int                       `json:"category_id" form:"category_id" validate:"required"`
	ProductImages   []ProductImagesRequest    `json:"product_images" form:"images"`
	ProductVideos   []ProductVideosRequest    `json:"product_videos" form:"videos"`
	ProductVariants *[]ProductVariantsRequest `json:"product_variants" form:"product_variants"`
}

type ProductPricingRequest struct {
	OriginalPrice   int  `json:"original_price" form:"original_price" validate:"required"`
	DiscountPercent *int `json:"discount_percent" form:"discount_percent"`
}

type ProductVariantsRequest struct {
	Stock int    `json:"stock" form:"stock" validate:"required"`
	Size  string `json:"size" form:"size" validate:"required"`
}

type ProductImagesRequest struct {
	ImageUrl *string `json:"image_url" form:"image_url"`
}

type ProductVideosRequest struct {
	VideoUrl *string `json:"video_url" form:"video_url"`
}

type ProductResponse struct {
	ID              uuid.UUID                 `json:"id"`
	Name            string                    `json:"name"`
	Description     string                    `json:"description"`
	MinOrder        int                       `json:"min_order"`
	AuthorID        uuid.UUID                 `json:"author_id"`
	CategoryName    string                    `json:"category_name"` // Pastikan nama field benar
	ProductPricing  ProductPricingResponse    `json:"product_pricing"`
	ProductVariants *[]ProductVariantsResponse `json:"product_variants"`
	ProductImages   []ProductImagesResponse   `json:"product_images"`
	ProductVideos   []ProductVideosResponse   `json:"product_videos"`
}

type ProductPricingResponse struct {
	OriginalPrice   int      `json:"original_price"`
	DiscountPercent *int     `json:"discount_percent"`
	DiscountPrice   *float64 `json:"discount_price"`
}

type ProductVariantsResponse struct {
	Price int    `json:"price"`
	Stock int    `json:"stock"`
	Size  string `json:"size"`
}

type ProductImagesResponse struct {
	ImageUrl *string `json:"image_url"`
}

type ProductVideosResponse struct {
	VideoUrl *string `json:"video_url"`
}

type CategoryRequest struct {
	Name string `json:"name" form:"name" validate:"required"`
}

type CategoryResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
