package dto

import "github.com/google/uuid"

type EventResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type EventDetailResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}
