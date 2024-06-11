package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Articles struct {
	ID            uuid.UUID          `gorm:"primaryKey;type:uuid"`
	Title         string             `gorm:"type:varchar(100);not null"`
	Image         string             `gorm:"type:varchar(255)"`
	Content       string             `gorm:"type:text"`
	Tags          string             `gorm:"type:varchar(100)"`
	LikesCount    int                `gorm:"type:int"`
	CommentsCount int                `gorm:"type:int"`
	AuthorID      uuid.UUID          `gorm:"type:uuid;not null"`
	Comments      *[]ArticleComments `gorm:"foreignKey:ArticleID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type ArticleComments struct {
	ID              uuid.UUID `gorm:"primaryKey;type:uuid"`
	UserID          uuid.UUID `gorm:"type:uuid;not null"`
	ArticleID       uuid.UUID `gorm:"type:uuid;not null"`
	ParentCommentID *uuid.UUID
	Content         string             `gorm:"type:text"`
	Replies         *[]ArticleComments `gorm:"foreignKey:ParentCommentID"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

type ArticleLikes struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	ArticleID uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt time.Time
}