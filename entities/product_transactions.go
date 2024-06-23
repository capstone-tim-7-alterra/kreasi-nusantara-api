package entities

import (
	"time"

	"github.com/google/uuid"
)

type ProductTransaction struct {
	ID                string  `gorm:"primary_key;type:string"`
	CartId            uuid.UUID `gorm:"type:uuid;not null"`
	UserId            uuid.UUID `gorm:"type:uuid;not null"`
	TracsactionDate   time.Time
	TotalAmount       float64
	TransactionStatus string
	TransactionMethod string
	SnapURL           string
}

type UpdateTransaction struct {
	ID                string `gorm:"primary_key;type:uuid"`
	TransactionStatus string
	TransactionMethod string
}



