package productsadmin

import "mime/multipart"

type ProductRequest struct {
	ProductName string                `json:"product_name" form:"product_name" validate:"required"`
	Deskription string                `json:"deskription" form:"deskription" validate:"required"`
	Price       int                   `json:"price" form:"price" validate:"required"`
	Stock       int                   `json:"stock" form:"stock" validate:"required"`
	AuthorID    int                   `json:"author_id" form:"author_id" validate:"required"`
	CategoryID  int                   `json:"category_id" form:"category_id" validate:"required"`
	Image       *multipart.FileHeader ` form:"image" `
	Video       *multipart.FileHeader ` form:"video" `
}

type ProductResponse struct {
	ID          uint   `json:"id"`
	ProductName string `json:"product_name"`
	Deskription string `json:"deskription"`
	Price       int    `json:"price"`
	Stock       int    `json:"stock"`
	Image       string `json:"image"`
	Video       string `json:"video"`
	AuthorID    int    `json:"author_id"`
	LikesCount  int    `json:"likes_count"`
	CategoryID  int    `json:"category_id"`
}

type CategoryRequest struct {
	Name string `json:"name" form:"name" validate:"required"`
}

type CategoryResponse struct {
	ID   int   `json:"id"`
	Name string `json:"name"`
}
