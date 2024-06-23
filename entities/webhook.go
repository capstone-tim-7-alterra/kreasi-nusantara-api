package entities

import ()

type PaymentNotification struct {
	TransactionTime   string 
	TransactionStatus string 
	TransactionID     string 
	SignatureKey      string 
	PaymentType       string
	OrderID           string 
	MerchantID        string
	GrossAmount       string 
	FraudStatus       string 
	Currency          string 
	SettlementTime    string 
}