package admin

import (
	// "kreasi-nusantara-api/controllers"

	// "kreasi-nusantara-api/repositories"
	// "kreasi-nusantara-api/usecases"
	// "kreasi-nusantara-api/utils/password"
	"kreasi-nusantara-api/middlewares"
	"kreasi-nusantara-api/utils/token"
	"kreasi-nusantara-api/utils/validation"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitAdminRoute(g *echo.Group, db *gorm.DB, v *validation.Validator) {
	// userRepo := repositories.NewUserRepository(db)
	// passwordUtil := password.NewPasswordUtil()

	// userUseCase := usecases.NewUserUseCase(userRepo, passwordUtil)
	// userController := controllers.NewUserController(userUseCase, v)

	// Protected routes
	g.Use(echojwt.WithConfig(token.GetJWTConfig()), middlewares.IsAdmin)
	g.GET("/tes", func(c echo.Context) error {
		return c.JSON(200, "tes")
	})
}