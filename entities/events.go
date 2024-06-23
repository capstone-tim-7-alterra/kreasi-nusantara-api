package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Events struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name        string    `gorm:"type:varchar(100);not null"`
	CategoryID  int       `gorm:"type:int;not null"`
	LocationID  uuid.UUID `gorm:"type:uuid;not null"`
	Status      bool      `gorm:"default:true"`
	Date        time.Time
	Photos      []EventPhotos `gorm:"foreignKey:EventID;references:ID"`
	Description string        `gorm:"type:text;not null"`
	Prices      []EventPrices `gorm:"foreignKey:EventID;references:ID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Category    EventCategories
	Location    EventLocations
}

type EventCategories struct {
	ID        int      `gorm:"primaryKey;autoIncrement"`
	Name      string   `gorm:"type:varchar(100);not null"`
	Events    []Events `gorm:"foreignKey:CategoryID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type EventLocations struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid"`
	Building    string    `gorm:"type:varchar(100);not null"`
	Address     string    `gorm:"type:varchar(100);not null"`
	Province    string    `gorm:"type:varchar(100);not null"`
	City        string    `gorm:"type:varchar(100);not null"`
	Subdistrict string    `gorm:"type:varchar(100);not null"`
	PostalCode  string    `gorm:"type:varchar(100);not null"`
	Events      []Events  `gorm:"foreignKey:LocationID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type EventPhotos struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	EventID   uuid.UUID `gorm:"type:uuid;not null"`
	Image     *string    `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type EventPrices struct {
	ID           uuid.UUID       `gorm:"primaryKey;type:uuid"`
	EventID      uuid.UUID       `gorm:"type:uuid;not null"`
	TicketTypeID int             `gorm:"type:int;not null"`
	Price        int             `gorm:"type:int;not null"`
	NoOfTicket   int             `gorm:"type:int;not null"`
	TicketType   EventTicketType `gorm:"foreignKey:TicketTypeID"`
	Publish      time.Time
	EndPublish   time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type EventTicketType struct {
	ID        int           `gorm:"primaryKey;autoIncrement"`
	Name      string        `gorm:"type:varchar(100);not null"`
	Prices    []EventPrices `gorm:"foreignKey:TicketTypeID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
