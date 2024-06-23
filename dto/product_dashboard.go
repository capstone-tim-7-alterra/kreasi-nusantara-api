package dto

import "time"

type ProductDashboard struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Income        float64   `json:"income"`
	PaymentMethod string    `json:"payment_method"`
	Image         string    `json:"image"`
	Status        string    `json:"status"`
	Date          time.Time `json:"date"`
}

type ProductHeader struct {
	ProductSold   int     `json:"product_sold"`
	ProductProfit float64 `json:"income"`
	TicketSold    int     `json:"ticket_sold"`
	TicketProfit  float64 `json:"ticket_profit"`
}



type ProductChart struct {
	Name  string         `json:"name"`
	Value []ProductValue `json:"value"`
}

type ProductValue struct {
	Income float64 `json:"income"`
	Date   string  `json:"date"`
}
