package admin

import (
	// "kreasi-nusantara-api/controllers"
	"kreasi-nusantara-api/middlewares"
	// "kreasi-nusantara-api/repositories"
	// "kreasi-nusantara-api/usecases"
	// "kreasi-nusantara-api/utils/password"
	"kreasi-nusantara-api/utils/validation"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitUserRoute(g *echo.Group, db *gorm.DB, v *validation.Validator) {
	// userRepo := repositories.NewUserRepository(db)
	// passwordUtil := password.NewPasswordUtil()

	// userUseCase := usecases.NewUserUseCase(userRepo, passwordUtil)
	// userController := controllers.NewUserController(userUseCase, v)

	// Protected routes
	g.Use(middlewares.Auth, middlewares.IsAdmin)
	// g.POST("/register", adminController.Register)
}