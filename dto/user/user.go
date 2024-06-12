package dto

import (
	"mime/multipart"
)

type RegisterRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=32"`
}

type RegisterResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type VerifyOTPRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type LoginResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Email              string `json:"email" validate:"required,email"`
	NewPassword        string `json:"new_password" validate:"required,min=8,max=32"`
	ConfirmNewPassword string `json:"confirm_new_password" validate:"required,min=8,max=32"`
}

type UserProfileResponse struct {
	Username  string  `json:"username"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Email     string  `json:"email"`
	Phone     *string `json:"phone"`
	Photo     *string `json:"photo"`
	Bio       *string `json:"bio"`
}

type UpdateProfileRequest struct {
	FirstName string  `json:"first_name" validate:"required"`
	LastName  string  `json:"last_name" validate:"required"`
	Phone     *string `json:"phone"`
	Bio       *string `json:"bio"`
}

type UserProfilePhotoRequest struct {
	Photo *multipart.FileHeader `json:"photo" form:"photo"`
}

type ChangePasswordRequest struct {
	OldPassword        string `json:"old_password" validate:"required,min=8,max=32"`
	NewPassword        string `json:"new_password" validate:"required,min=8,max=32"`
	ConfirmNewPassword string `json:"confirm_new_password" validate:"required,min=8,max=32"`
}
