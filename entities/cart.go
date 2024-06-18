package entities

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID        uuid.UUID    `gorm:"primary_key;type:uuid"`
	UserID    uuid.UUID    `gorm:"type:uuid;not null"`
	Items     []CartItems `gorm:"foreignKey:CartID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CartItems struct {
	ID               uuid.UUID `gorm:"primary_key;type:uuid"`
	CartID           uuid.UUID `gorm:"type:uuid;not null"`
	ProductVariantID uuid.UUID `gorm:"type:uuid;not null"`
	Quantity         int       `gorm:"type:int;not null"`
	ProductVariant   ProductVariants
}
