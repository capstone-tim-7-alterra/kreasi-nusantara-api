package controllers

import (
	"kreasi-nusantara-api/entities"
	"kreasi-nusantara-api/usecases"

	"github.com/labstack/echo/v4"
)

type WebhookController struct {
	WebhookUsecase usecases.WebhookUsecase
}

func NewWebhookController(webhookUsecase usecases.WebhookUsecase) *WebhookController {
	return &WebhookController{
		WebhookUsecase: webhookUsecase,
	}
}

func (w *WebhookController) HandleNotification() echo.HandlerFunc {
	return func(c echo.Context) error {
		var notification entities.PaymentNotification
		err := c.Bind(&notification)
		if err != nil {
			return echo.NewHTTPError(400, err.Error())
		}
		err = w.WebhookUsecase.HandleNotification(c, notification)
		if err != nil {
			return echo.NewHTTPError(500, err.Error())
		}
		return c.JSON(200, map[string]string{
			"message": "success",
		})
	}
}
