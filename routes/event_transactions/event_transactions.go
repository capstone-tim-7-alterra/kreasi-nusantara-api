package event_transactions

import (
	"kreasi-nusantara-api/config"
	"kreasi-nusantara-api/controllers"
	"kreasi-nusantara-api/drivers/redis"
	"kreasi-nusantara-api/repositories"
	"kreasi-nusantara-api/usecases"
	"kreasi-nusantara-api/utils/token"
	"kreasi-nusantara-api/utils/validation"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitEventTransactionsRoute(g *echo.Group, db *gorm.DB, v *validation.Validator) {
	tokenUtil := token.NewTokenUtil()
	config := config.InitConfigMidtrans()
	redisClient := redis.NewRedisClient()

	eventAdminRepository := repositories.NewEventAdminRepository(db)
	eventTransactionRepo := repositories.NewEventTransactionRepository(db)
	eventTransactionUseCase := usecases.NewEventTransactionUseCase(eventTransactionRepo, eventAdminRepository, *redisClient,config)

	eventTransactionController := controllers.NewEventTransactionController(eventTransactionUseCase, v, tokenUtil)

	g.Use(echojwt.WithConfig(token.GetJWTConfig()))
	g.POST("/event-transactions", eventTransactionController.CreateEventTransaction)
	g.GET("/event-transactions/:id", eventTransactionController.GetEventTransactionById)
}
