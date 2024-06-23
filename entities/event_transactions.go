package entities

import (
	"time"

	"github.com/google/uuid"
)

type EventTransaction struct {
	ID                uuid.UUID `gorm:"primary_key;type:uuid"`
	EventPriceID      uuid.UUID `gorm:"type:uuid;not null"`
	UserId            uuid.UUID `gorm:"type:uuid;not null"`
	TransactionDate   time.Time
	TotalAmount       float64
	TransactionStatus string
	TransactionMethod string
	SnapURL           string
	Buyer             EventTransactionBuyer
}

type EventTransactionBuyer struct {
	ID                 uuid.UUID `gorm:"primary_key;type:uuid"`
	EventTransactionID uuid.UUID `gorm:"type:uuid;not null"`
	IdentityNumber     string    `gorm:"type:varchar(100);not null"`
	FullName           string    `gorm:"type:varchar(100);not null"`
	Email              string    `gorm:"type:varchar(100);not null"`
	Phone              string    `gorm:"type:varchar(100);not null"`
}
