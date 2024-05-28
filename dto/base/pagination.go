package dto

import "time"

type PaginationRequest struct {
	Limit int        `json:"limit" validate:"required,gt=0"`
	Page  int        `json:"offset" validate:"required,gt=0"`
	From  *time.Time `json:"from"`
	To    *time.Time `json:"to"`
}

type PaginationResponse struct {
	BaseResponse
	Pagination *PaginationMetadata `json:"pagination"`
	Link       *Link               `json:"link"`
}

type Link struct {
	Next string `json:"next,omitempty"`
	Prev string `json:"prev,omitempty"`
}

type PaginationMetadata struct {
	CurrentPage int   `json:"current_page"`
	TotalPage   int   `json:"total_page"`
	TotalData   int64 `json:"total_data"`
}