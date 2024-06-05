package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Products struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid"`
	ProductName string    `gorm:"type:varchar(100);not null"`
	Description string    `gorm:"type:varchar(255);not null"`
	Price       int       `gorm:"type:int;not null"`
	Stock       int       `gorm:"type:int;not null"`
	Image       *string   `gorm:"type:varchar(255)"`
	Video       *string   `gorm:"type:varchar(255)"`
	LikesCount  *int      `gorm:"type:int"`
	AuthorID    uuid.UUID `gorm:"type:uuid;not null"`
	CategoryID  int       `gorm:"type:int;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type Category struct {
	ID        int        `gorm:"primaryKey;autoIncrement"`
	Name      string     `gorm:"type:varchar(100);not null"`
	Products  []Products `gorm:"foreignKey:CategoryID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
