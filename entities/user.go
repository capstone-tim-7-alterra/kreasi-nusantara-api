package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID              uuid.UUID          `gorm:"primaryKey;type:uuid"`
	Username        string             `gorm:"type:varchar(64);unique;not null"`
	FirstName       string             `gorm:"type:varchar(100); not null"`
	LastName        string             `gorm:"type:varchar(100); not null"`
	Email           string             `gorm:"type:varchar(100);unique;not null"`
	Password        string             `gorm:"type:varchar(100); not null"`
	Phone           *string            `gorm:"type:varchar(20)"`
	Photo           *string            `gorm:"type:varchar(255)"`
	Gender          *string            `gorm:"type:char(1)"`
	DateOfBirth     *time.Time         `gorm:"type:date"`
	IsVerified      bool               `gorm:"default:false"`
	Addresses       *[]UserAddresses   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ArticleComments *[]ArticleComments `gorm:"foreignKey:ArticleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}
