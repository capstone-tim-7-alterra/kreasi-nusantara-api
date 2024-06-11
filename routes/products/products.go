package products

import (
	"kreasi-nusantara-api/controllers"
	"kreasi-nusantara-api/repositories"
	"kreasi-nusantara-api/usecases"
	"kreasi-nusantara-api/utils/token"
	"kreasi-nusantara-api/utils/validation"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitProductsRoute(g *echo.Group, db *gorm.DB, v *validation.Validator) {
	productRepo := repositories.NewProductRepository(db)
	productUseCase := usecases.NewProductUseCase(productRepo)
	productController := controllers.NewProductController(productUseCase, v)

	g.Use(echojwt.WithConfig(token.GetJWTConfig()))
	g.GET("/products", productController.GetProducts)
	g.GET("/products/:product_id", productController.GetProductByID)
	g.GET("/products/search", productController.SearchProducts)
	g.GET("/products/category/:category_id", productController.GetProductsByCategory)
}