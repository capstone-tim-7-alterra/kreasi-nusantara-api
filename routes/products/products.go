package products

import (
	"kreasi-nusantara-api/controllers"
	"kreasi-nusantara-api/drivers/redis"
	"kreasi-nusantara-api/repositories"
	"kreasi-nusantara-api/usecases"
	"kreasi-nusantara-api/utils/token"
	"kreasi-nusantara-api/utils/validation"
	"os"
	oai "kreasi-nusantara-api/drivers/openai"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitProductsRoute(g *echo.Group, db *gorm.DB, v *validation.Validator) {
	oaiKey := os.Getenv("OPENAI_API_KEY")

	oaiService := oai.NewOpenAIClient(oaiKey)
	redisClient := redis.NewRedisClient()

	productRepo := repositories.NewProductRepository(db)
	productUseCase := usecases.NewProductUseCase(productRepo)

	tokenUtil := token.NewTokenUtil()

	cartRepo := repositories.NewCartRepository(db)

	productController := controllers.NewProductController(productUseCase, v, tokenUtil)

	recommendationUseCase := usecases.NewRecommendationUseCase(oaiService, *redisClient, productRepo, cartRepo)
	recommendationController := controllers.NewRecommendationController(recommendationUseCase, tokenUtil)

	g.Use(echojwt.WithConfig(token.GetJWTConfig()))
	g.GET("/products", productController.GetProducts)
	g.GET("/products/:product_id", productController.GetProductByID)
	g.GET("/products/search", productController.SearchProducts)
	g.GET("/products/category/:category_id", productController.GetProductsByCategory)

	g.POST("/products/:product_id/reviews", productController.CreateProductReview)
	g.GET("/products/:product_id/reviews", productController.GetProductReviews)
	g.GET("/products/recommendation", recommendationController.GetProductRecommendation)
}