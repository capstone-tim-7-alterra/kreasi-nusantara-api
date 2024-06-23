package usecases

import (

	"kreasi-nusantara-api/entities"
	"kreasi-nusantara-api/repositories"

	"github.com/labstack/echo/v4"
)

type WebhookUsecase interface {
	HandleNotification(c echo.Context, webhook entities.PaymentNotification) error
}

type webhookUsecase struct {
	webhookRepository repositories.WebhookRepository
}

func NewWebhookUsecase(webhookRepository repositories.WebhookRepository) WebhookUsecase {
	return &webhookUsecase{
		webhookRepository: webhookRepository,
	}
}

func (u *webhookUsecase) HandleNotification(c echo.Context, webhook entities.PaymentNotification) error {
	transactionStatus := webhook.TransactionStatus
	fraudStatus := webhook.FraudStatus
	transactionData := entities.ProductTransaction{
		ID:                webhook.OrderID,
		TransactionMethod: webhook.PaymentType,
	}

	if transactionStatus == "capture" {
		if fraudStatus == "accept" {
			transactionData.TransactionStatus = "paid"
		} else if fraudStatus == "challenge" {
			transactionData.TransactionStatus = "challenge"
		} else if fraudStatus == "reject" {
			transactionData.TransactionStatus = "rejected"
		}
	} else if transactionStatus == "settlement" {
		transactionData.TransactionStatus = "paid"
	} else if transactionStatus == "deny" {
		transactionData.TransactionStatus = "rejected"
	} else if transactionStatus == "cancel" || transactionStatus == "expire" {
		transactionData.TransactionStatus = "canceled"
	} else if transactionStatus == "pending" {
		transactionData.TransactionStatus = "pending"
	}
	return u.webhookRepository.HandleNotification(c.Request().Context(), webhook, transactionData)
}