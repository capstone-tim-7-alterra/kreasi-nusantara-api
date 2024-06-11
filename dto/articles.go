package dto

import (
	"kreasi-nusantara-api/entities"
	"time"

	"github.com/google/uuid"
)

type ArticleResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
}

type ArticleDetailResponse struct {
	ID            uuid.UUID      `json:"id"`
	Title         string         `json:"title"`
	Content       string         `json:"content"`
	LikesCount    int            `json:"likes_count"`
	CommentsCount int            `json:"comments_count"`
	CreatedAt     time.Time      `json:"created_at"`
	Author        entities.Admin `json:"author"`
}

type ArticleCommentResponse struct {
	ID        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type ArticleCommentRequest struct {
	Content string `json:"content"`
}