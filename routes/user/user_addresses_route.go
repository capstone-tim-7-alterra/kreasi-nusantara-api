package user

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

func InitUserAddressesRoute(g *echo.Group, db *gorm.DB, v *validation.Validator) {
	userAddressRepo := repositories.NewUserAddressRepository(db)
	userAddressUseCase := usecases.NewUserAddressUseCase(userAddressRepo)
	tokenUtil := token.NewTokenUtil()

	userAddressController := controllers.NewUserAddressController(userAddressUseCase, v, tokenUtil)

	// Protected routes
	g.Use(echojwt.WithConfig(token.GetJWTConfig()))
	g.GET("/users/me/addresses", userAddressController.GetUserAddresses)
	g.GET("/users/me/addresses/:address_id", userAddressController.GetUserAddressByID)
	g.POST("/users/me/addresses", userAddressController.CreateUserAddress)
	g.PUT("/users/me/addresses/:address_id", userAddressController.UpdateUserAddress)
	g.DELETE("/users/me/addresses/:address_id", userAddressController.DeleteUserAddress)
}