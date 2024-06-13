package dto

import (

	"github.com/google/uuid"
)

type EventResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type EventDetailResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type EventCreateRequest struct {
	Name        string               `json:"name" form:"name" validate:"required"`
	Description string               `json:"description" form:"description" validate:"required"`
	CategoryID  int                  `json:"category_id" form:category_id" validate:"required"`
	LocationID  uuid.UUID            `json:"location_id" form:"location_id" validate:"required"`
	Date        string               `json:"date" form:"date" validate:"required"`
	Photos      []EventPhotosRequest `json:"photos" form:"photos"`
	Prices      []EventPricesRequest `json:"prices" form:"prices"`
}

type EventPhotosRequest struct {
	Image string  `form:"image"`
}

type EventPricesRequest struct {
	TicketTypeID int `json:"ticket_type_id" form:"ticket_type_id" validate:"required"`
	Price        int `json:"price" form:"price" validate:"required"`
	NoOfTicket   int `json:"no_of_ticket" form:"no_of_ticket" validate:"required"`
	Publish      string `json:"publish" form:"publish" validate:"required"`
	EndPublish   string `json:"end_publish" form:"end_publish" validate:"required"`
}

type EventTicketTypeRequest struct {
	Name string `json:"name" form:"name" validate:"required"`
}

type EventLocationsRequest struct {
	Building    string `json:"building" form:"building" validate:"required"`
	Address     string `json:"address" form:"address" validate:"required"`
	City        string `json:"city" form:"city" validate:"required"`
	Subdistrict string `json:"subdistrict" form:"subdistrict" validate:"required"`
	PostalCode  string `json:"postal_code" form:"postal_code" validate:"required"`
}

type EventCategoriesRequest struct {
	Name string `json:"name" form:"name" validate:"required"`
}

type EventAdminResponse struct {
	ID       uuid.UUID             `json:"id"`
	Name     string                `json:"name"`
	Status   string                `json:"status"`
	Date     string                `json:"date"`
	Ticket   string                `json:"ticket"`
	Location string                `json:"location"`
	Photos   []EventPhotosResponse `json:"photos"`
}

type EventPhotosResponse struct {
	ID       uuid.UUID `json:"id"`
	ImageUrl string    `json:"image_url"`
}
