package dto

import (
	"mime/multipart"
)

type RegisterRequest struct {
	FirstName    string                `json:"first_name" form:"first_name" validate:"required"`
	LastName     string                `json:"last_name" form:"last_name" validate:"required"`
	Username     string                `json:"username" form:"username" validate:"required"`
	Email        string                `json:"email" form:"email" validate:"required,email"`
	Password     string                `json:"password" form:"password" validate:"required,min=8,max=32"`
	IsSuperAdmin bool                  `json:"is_super_admin" form:"is_super_admin"`
	Image        *multipart.FileHeader ` form:"image" `
}

type RegisterResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}

type LoginResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

type AdminResponse struct {
	ID           string  `json:"id"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Username     string  `json:"username"`
	Email        string  `json:"email"`
	IsSuperAdmin bool    `json:"is_super_admin"`
	Photo        *string `json:"photo"`
	CreatedAt    string  `json:"created_at"`
}

type AdminAvatarResponse struct {
	Photo *string `json:"photo"`
	Name  string  `json:"name`
}

type UpdateAdminRequest struct {
	FirstName    string                `json:"first_name" form:"first_name" validate:"required"`
	LastName     string                `json:"last_name" form:"last_name" validate:"required"`
	Username     string                `json:"username" form:"username" validate:"required"`
	Email        string                `json:"email" form:"email" validate:"required,email"`
	Password     string                `json:"password" form:"password" validate:"required,min=8,max=32"`
	IsSuperAdmin bool                  `json:"is_super_admin" form:"is_super_admin"`
	Photo        *multipart.FileHeader ` form:"photo" `
}
