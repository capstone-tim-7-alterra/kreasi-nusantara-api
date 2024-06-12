package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Admin struct {
	ID           uuid.UUID  `gorm:"primaryKey;type:uuid"`
	Username     string     `gorm:"type:varchar(64);unique;not null"`
	FirstName    string     `gorm:"type:varchar(100); not null"`
	LastName     string     `gorm:"type:varchar(100); not null"`
	Email        string     `gorm:"type:varchar(100);unique;not null"`
	Password     string     `gorm:"type:varchar(100); not null"`
	Photo        *string    `gorm:"type:varchar(255); not null"`
	Token        string     `gorm:"-"`
	IsSuperAdmin bool       `gorm:"default:false"`
	Products     []Products `gorm:"foreignKey:AuthorID"`
	Articles     []Articles `gorm:"foreignKey:AuthorID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
