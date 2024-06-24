package dto

import (
	"time"

	"github.com/google/uuid"
)

type ProductDashboard struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Income        float64   `json:"income"`
	PaymentMethod string    `json:"payment_method"`
	Image         string    `json:"image"`
	Status        string    `json:"status"`
	Date          time.Time `json:"date"`
}

type EventDashboard struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	Income        float64   `json:"income"`
	PaymentMethod string    `json:"payment_method"`
	Image         string    `json:"image"`
	Status        string    `json:"status"`
	Date          time.Time `json:"date"`
}

type ProductHeader struct {
	ProductSold   int     `json:"product_sold"`
	ProductProfit float64 `json:"product_profit"`
	TicketSold    int     `json:"ticket_sold"`
	TicketProfit  float64 `json:"ticket_profit"`
	TotalTicket   int     `json:"total_ticket"`
	DeletedTicket int     `json:"deleted_ticket"`
	TotalLikes    int     `json:"total_likes"`
	TotalComments int     `json:"total_comments"`
	TotalVisitors int     `json:"total_visitors"`
	TotalShares   int     `json:"total_shares"`
}

type ProductChart struct {
	Name  string         `json:"name"`
	Value []ProductValue `json:"value"`
}

type ProductValue struct {
	Income float64 `json:"income"`
	Date   string  `json:"date"`
}

type EventChart struct {
	Name  string       `json:"name"`
	Value []EventValue `json:"value"`
}

type EventValue struct {
	Income float64 `json:"income"`
	Date   string  `json:"date"`
}
