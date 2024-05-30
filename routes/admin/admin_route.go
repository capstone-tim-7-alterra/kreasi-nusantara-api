package admin

import (
	"kreasi-nusantara-api/config"
	"kreasi-nusantara-api/controllers"
	"kreasi-nusantara-api/drivers/cloudinary"
	// "kreasi-nusantara-api/middlewares"
	"kreasi-nusantara-api/repositories"
	"kreasi-nusantara-api/usecases"
	"kreasi-nusantara-api/utils/password"
	"kreasi-nusantara-api/utils/validation"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitUserRoute(g *echo.Group, db *gorm.DB, v *validation.Validator) {

	adminRepo := repositories.NewAdminRepository(db)
	passwordUtil := password.NewPasswordUtil()
	cloudinaryInstance, _ := config.SetupCloudinary()
	cloudinaryService := cloudinary.NewCloudinaryService(cloudinaryInstance)

	adminUseCase := usecases.NewAdminUsecase(adminRepo, passwordUtil, cloudinaryService)
	adminController := controllers.NewAdminController(adminUseCase, v)

	// Public routes
	g.POST("admin/register", adminController.Register)

	// Protected routes
	// g.Use(middlewares.Auth, middlewares.IsAdmin)
	// g.POST("/register", adminController.Register)
}
