package cart

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

func InitCartRoute(g *echo.Group, db *gorm.DB, v *validation.Validator) {
	cartRepo := repositories.NewCartRepository(db)

	cartUseCase := usecases.NewCartUseCase(cartRepo)
	tokenUtil := token.NewTokenUtil()

	cartController := controllers.NewCartController(cartUseCase, v, tokenUtil)

	g.Use(echojwt.WithConfig(token.GetJWTConfig()))
	g.POST("/carts/items", cartController.AddToCart)
	g.GET("/carts", cartController.GetCartItems)
	g.PUT("/carts/:cartItemID", cartController.UpdateCartItems)
	g.DELETE("/carts/:cartItemID", cartController.DeleteCartItem)
	g.GET("/carts-all", cartController.GetAllCarts)
}
