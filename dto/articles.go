package dto

import (
	"time"

	"github.com/google/uuid"
)

type ArticleResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
}

type ArticleAdminResponse struct {
	ID        uuid.UUID `json:"id"`
	Tags      string    `json:"tags"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	Image     string    `json:"image"`
	CreatedAt string    `json:"created_at"`
}

type ArticleDetailResponse struct {
	ID            uuid.UUID         `json:"id"`
	Title         string            `json:"title"`
	Content       string            `json:"content"`
	LikesCount    int               `json:"likes_count"`
	CommentsCount int               `json:"comments_count"`
	CreatedAt     time.Time         `json:"created_at"`
	Author        AuthorInformation `json:"author"`
}

type ArticleCommentResponse struct {
	ID        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type ArticleCommentRequest struct {
	Content string `json:"content"`
}

type ArticleRequest struct {
	Title   string `json:"title" form:"title"`
	Image   string `json:"image" form:"image"`
	Content string `json:"content" form:"content"`
	Tags    string `json:"tags" form:"tags"`
	Author  string `json:"author" form:"author"`
}

type AuthorInformation struct {
	ImageURL string `json:"image_url"`
	Username string `json:"username"`
}
