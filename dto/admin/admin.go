package dto

import "mime/multipart"

type RegisterRequest struct {
	FirstName    string `json:"first_name" form:"first_name" validate:"required"`
	LastName     string `json:"last_name" form:"last_name" validate:"required"`
	Username     string `json:"username" form:"username" validate:"required"`
	Email        string `json:"email" form:"email" validate:"required,email"`
	Password     string `json:"password" form:"password" validate:"required,min=8,max=32,alphanum"`
	IsSuperAdmin bool   `json:"is_super_admin" form:"is_super_admin"`
	Image        *multipart.FileHeader ` form:"image" `
}

type RegisterResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
