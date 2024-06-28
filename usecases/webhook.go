package usecases

import (
	"kreasi-nusantara-api/drivers/redis"
	"kreasi-nusantara-api/entities"
	"kreasi-nusantara-api/repositories"

	"github.com/labstack/echo/v4"
)

type WebhookUsecase interface {
	HandleNotification(c echo.Context, webhook entities.PaymentNotification) error
}

type webhookUsecase struct {
	webhookRepository repositories.WebhookRepository
	redisClient redis.RedisClient
}

func NewWebhookUsecase(webhookRepository repositories.WebhookRepository, redisClient redis.RedisClient) WebhookUsecase {
	return &webhookUsecase{
		webhookRepository: webhookRepository,
		redisClient:       redisClient,
	}
}

func (u *webhookUsecase) HandleNotification(c echo.Context, webhook entities.PaymentNotification) error {
	transactionStatus := webhook.TransactionStatus
	fraudStatus := webhook.FraudStatus

	key := "transaction-" + webhook.OrderID

	res, err := u.redisClient.Get(key)
	if err != nil && err.Error() != "redis: nil" {
		return err
	}

	transactionUpdate := entities.UpdateTransaction{
		ID:                webhook.OrderID,
		TransactionStatus: webhook.TransactionStatus,
		TransactionMethod: webhook.PaymentType,
	}

	if transactionStatus == "capture" {
		if fraudStatus == "accept" {
			transactionUpdate.TransactionStatus = "paid"
		} else if fraudStatus == "challenge" {
			transactionUpdate.TransactionStatus = "challenge"
		} else if fraudStatus == "reject" {
			transactionUpdate.TransactionStatus = "rejected"
		}
	} else if transactionStatus == "settlement" {
		transactionUpdate.TransactionStatus = "paid"
	} else if transactionStatus == "deny" {
		transactionUpdate.TransactionStatus = "rejected"
	} else if transactionStatus == "cancel" || transactionStatus == "expire" {
		transactionUpdate.TransactionStatus = "canceled"
	} else if transactionStatus == "pending" {
		transactionUpdate.TransactionStatus = "pending"
	}

	if res == "event" {
		return u.webhookRepository.HandleNotification(c.Request().Context(), webhook, transactionUpdate, "event_transactions")
	} else if res == "tes123" {
		return u.webhookRepository.HandleNotification(c.Request().Context(), webhook, transactionUpdate, "single_product_transactions")
	}
	return u.webhookRepository.HandleNotification(c.Request().Context(), webhook, transactionUpdate, "product_transactions")
}
