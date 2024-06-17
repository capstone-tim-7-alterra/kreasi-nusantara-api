package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Articles struct {
	ID            uuid.UUID                `gorm:"primaryKey;type:uuid"`
	Title         string                   `gorm:"type:varchar(100);not null"`
	Image         string                   `gorm:"type:varchar(255)"`
	Content       string                   `gorm:"type:text"`
	Tags          string                   `gorm:"type:varchar(100)"`
	LikesCount    int                      `gorm:"type:int"`
	CommentsCount int                      `gorm:"type:int"`
	AuthorID      uuid.UUID                `gorm:"type:uuid;not null"`
	Comments      *[]ArticleComments       `gorm:"foreignKey:ArticleID"`
	Replies       *[]ArticleCommentReplies `gorm:"foreignKey:ArticleID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	Author        *Admin
}

type ArticleComments struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	ArticleID uuid.UUID `gorm:"type:uuid;not null"`
	Content   string    `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt           `gorm:"index"`
	Replies   *[]ArticleCommentReplies `gorm:"foreignKey:CommentID"`
}

type ArticleCommentReplies struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	ArticleID uuid.UUID `gorm:"type:uuid;not null"`
	CommentID uuid.UUID `gorm:"type:uuid;not null"`
	Content   string    `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type ArticleLikes struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	ArticleID uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt time.Time
}
