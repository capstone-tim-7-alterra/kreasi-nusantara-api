package user

import (
	"kreasi-nusantara-api/controllers"
	// "kreasi-nusantara-api/middlewares"
	"kreasi-nusantara-api/repositories"
	"kreasi-nusantara-api/usecases"
	"kreasi-nusantara-api/utils/password"
	"kreasi-nusantara-api/utils/validation"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitUserRoute(g *echo.Group, db *gorm.DB, v *validation.Validator) {
	userRepo := repositories.NewUserRepository(db)
	passwordUtil := password.NewPasswordUtil()

	userUseCase := usecases.NewUserUseCase(userRepo, passwordUtil)
	userController := controllers.NewUserController(userUseCase, v)

	// Public routes
	g.POST("/register", userController.Register)

	// Protected routes
	// g.Use(middlewares.Auth)
	// g.GET("/profile", userController.GetProfile)
}
