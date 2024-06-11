package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Products struct {
	ID              uuid.UUID          `gorm:"primaryKey;type:uuid"`
	Name            string             `gorm:"type:varchar(100)"`
	Description     string             `gorm:"type:varchar(255)"`
	MinOrder        int                `gorm:"type:int"`
	AuthorID        uuid.UUID          `gorm:"type:uuid"`
	CategoryID      int                `gorm:"type:int"`
	ProductPricing  ProductPricing     `gorm:"foreignKey:ProductID;references:ID"`
	ProductVariants *[]ProductVariants `gorm:"foreignKey:ProductID;references:ID"`
	ProductImages   []ProductImages    `gorm:"foreignKey:ProductID;references:ID"`
	ProductVideos   []ProductVideos    `gorm:"foreignKey:ProductID;references:ID"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

type ProductPricing struct {
	ID              uuid.UUID `gorm:"primaryKey;type:uuid"`
	ProductID       uuid.UUID `gorm:"type:uuid"`
	OriginalPrice   int       `gorm:"type:int"`
	DiscountPercent *int      `gorm:"type:int"`
	DiscountPrice   *float64  `gorm:"type:decimal(10,2)"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

type ProductVariants struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	ProductID uuid.UUID `gorm:"type:uuid"`
	Price     int      `gorm:"type:int"`
	Stock     int      `gorm:"type:int"`
	Size      string   `gorm:"type:varchar(255)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type ProductImages struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	ProductID uuid.UUID `gorm:"type:uuid"`
	ImageUrl  *string   `gorm:"type:varchar(255)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type ProductVideos struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	ProductID uuid.UUID `gorm:"type:uuid"`
	VideoUrl  *string   `gorm:"type:varchar(255)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type ProductCategory struct {
	ID        int        `gorm:"primaryKey;autoIncrement"`
	Name      string     `gorm:"type:varchar(100);not null"`
	Products  []Products `gorm:"foreignKey:CategoryID;references:ID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
