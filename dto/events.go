package dto

import (
	"github.com/google/uuid"
)

type EventResponse struct {
	ID       uuid.UUID           `json:"id"`
	Name     string              `json:"name"`
	Image    string              `json:"image"`
	Category string              `json:"category"`
	Location EventLocationDetail `json:"location"`
	Date     string              `json:"date"`
	MinPrice int                 `json:"min_price"`
}

type EventDetailResponse struct {
	ID          uuid.UUID             `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Images      []string              `json:"images"`
	Location    EventLocationDetail   `json:"location"`
	Date        string                `json:"date"`
	Ticket      []EventPricesResponse `json:"ticket"`
}

type EventLocationDetail struct {
	Building    string `json:"building"`
	Subdistrict string `json:"subdistrict"`
	City        string `json:"city"`
}

type EventRequest struct {
	Name        string               `json:"name" form:"name" validate:"required"`
	Description string               `json:"description" form:"description" validate:"required"`
	CategoryID  int                  `json:"category_id" form:"category_id" validate:"required"`
	Date        string               `json:"date" form:"date" validate:"required"`
	Prices      []EventPricesRequest `json:"prices" form:"prices" validate:"required"`
	Photos      []EventPhotosRequest `json:"photos" form:"photos" validate:"required"`
	Location    EventLocationRequest `json:"location" form:"location" validate:"required"`
}

type EventLocationRequest struct {
	Building    string `json:"building" form:"building" validate:"required"`
	Address     string `json:"address" form:"address" validate:"required"`
	Province    string `json:"province" form:"province" validate:"required"`
	City        string `json:"city" form:"city" validate:"required"`
	Subdistrict string `json:"subdistrict" form:"subdistrict" validate:"required"`
	PostalCode  string `json:"postal_code" form:"postal_code" validate:"required"`
}

type EventCategoryRequest struct {
	Name string `json:"name" form:"name" validate:"required"`
}

type EventPhotosRequest struct {
	Image string `form:"image"`
}

type EventPricesRequest struct {
	Price        int    `json:"price" form:"price" validate:"required"`
	TicketTypeID int    `json:"ticket_type_id" form:"ticket_type_id" validate:"required"`
	NoOfTicket   int    `json:"no_of_ticket" form:"no_of_ticket" validate:"required"`
	Publish      string `json:"publish" form:"publish" validate:"required"`
	EndPublish   string `json:"end_publish" form:"end_publish" validate:"required"`
}

type EventTicketTypeRequest struct {
	Name string `json:"name" form:"name" validate:"required"`
}

type EventCategoriesRequest struct {
	Name string `json:"name" form:"name" validate:"required"`
}

type EventAdminResponse struct {
	ID         uuid.UUID             `json:"id"`
	Name       string                `json:"name"`
	Status     string                `json:"status"`
	TypeTicket string                `json:"type_ticket"`
	Date       string                `json:"date"`
	Location   string                `json:"location"`
	Photos     []EventPhotosResponse `json:"photos"`
}

type EventAdminDetailResponse struct {
	ID          uuid.UUID               `json:"id"`
	Name        string                  `json:"name"`
	Status      string                  `json:"status"`
	Date        string                  `json:"date"`
	Description string                  `json:"description"`
	Category    EventCategoriesResponse `json:"category"`
	Ticket      []EventPricesResponse   `json:"ticket"`
	Location    EventLocationResponse   `json:"location"`
	Photos      []EventPhotosResponse   `json:"photos"`
}

type EventLocationResponse struct {
	ID          uuid.UUID `json:"id"`
	Building    string    `json:"building"`
	Address     string    `json:"address" `
	Province    string    `json:"province" `
	City        string    `json:"city" `
	Subdistrict string    `json:"subdistrict" `
	PostalCode  string    `json:"postal_code" `
}

type EventPricesResponse struct {
	ID         uuid.UUID               `json:"id"`
	Price      int                     `json:"price" `
	TicketType EventTicketTypeResponse `json:"ticket_type" `
	NoOfTicket int                     `json:"no_of_ticket" `
	Publish    string                  `json:"publish" `
	EndPublish string                  `json:"end_publish" `
}

type EventPhotosResponse struct {
	ID       uuid.UUID `json:"id"`
	ImageUrl string    `json:"image_url"`
}

type EventTicketTypeResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type EventCategoriesResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Province struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProvinceResponse struct {
	Code     string     `json:"code"`
	Messages string     `json:"messages"`
	Value    []Province `json:"value"`
}

type District struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DistrictResponse struct {
	Code     string     `json:"code"`
	Messages string     `json:"messages"`
	Value    []District `json:"value"`
}

type Subdistrict struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SubdistrictResponse struct {
	Code     string        `json:"code"`
	Messages string        `json:"messages"`
	Value    []Subdistrict `json:"value"`
}
