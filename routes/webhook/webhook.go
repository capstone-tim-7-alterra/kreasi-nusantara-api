package webhook

import (
	"kreasi-nusantara-api/controllers"
	"kreasi-nusantara-api/drivers/redis"
	"kreasi-nusantara-api/repositories"
	"kreasi-nusantara-api/usecases"

	// "kreasi-nusantara-api/utils/token"

	// echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitWebhookRoute(g *echo.Group, db *gorm.DB) {
	redisClient := redis.NewRedisClient()

	webhookRepo := repositories.NewWebhookRepository(db)
	webhookUsecase := usecases.NewWebhookUsecase(webhookRepo, *redisClient)
	webhookController := controllers.NewWebhookController(webhookUsecase)

	// g.Use(echojwt.WithConfig(token.GetJWTConfig()))
	g.POST("/midtrans-notification", webhookController.HandleNotification())
}
