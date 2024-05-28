package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserAddresses struct {
	ID            uuid.UUID `gorm:"primaryKey;type:uuid"`
	UserID        uuid.UUID `gorm:"type:uuid;not null"`
	Label         string    `gorm:"type:varchar(10);not null"`
	RecipientName string    `gorm:"type:varchar(100);not null"`
	Phone         string    `gorm:"type:varchar(20);not null"`
	Address       string    `gorm:"type:varchar(255);not null"`
	City          string    `gorm:"type:varchar(100);not null"`
	Province      string    `gorm:"type:varchar(100);not null"`
	PostalCode    string    `gorm:"type:varchar(20);not null"`
	IsPrimary     bool      `gorm:"default:false"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}