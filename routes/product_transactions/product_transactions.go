package product_transactions

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

func InitProductTransactionsRoute(g *echo.Group, db *gorm.DB, v *validation.Validator) {
	tokenUtil := token.NewTokenUtil()
	cartRepo := repositories.NewCartRepository(db)
	cartUseCase := usecases.NewCartUseCase(cartRepo)
	config := config.InitConfigMidtrans()
	redisClient := redis.NewRedisClient()

	productTransactionRepo := repositories.NewProductTransactionRepository(db)
	productTransactionUseCase := usecases.NewProductTransactionUseCase(productTransactionRepo, cartUseCase, tokenUtil, *redisClient ,config)
	productTransactionController := controllers.NewProductTransactionController(productTransactionUseCase, v)

	g.Use(echojwt.WithConfig(token.GetJWTConfig()))
	g.POST("/product-transactions", productTransactionController.CreateProductTransaction)
	g.GET("/product-transactions/:id", productTransactionController.GetProductTransactionById)

}
