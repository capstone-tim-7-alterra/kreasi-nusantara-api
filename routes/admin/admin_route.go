package admin

import (
	"kreasi-nusantara-api/config"
	"kreasi-nusantara-api/controllers"
	"kreasi-nusantara-api/drivers/cloudinary"
	"kreasi-nusantara-api/middlewares"

	// "kreasi-nusantara-api/middlewares"
	"kreasi-nusantara-api/repositories"
	"kreasi-nusantara-api/usecases"
	"kreasi-nusantara-api/utils/password"
	"kreasi-nusantara-api/utils/token"
	"kreasi-nusantara-api/utils/validation"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitAdminRoute(g *echo.Group, db *gorm.DB, v *validation.Validator) {

	adminRepo := repositories.NewAdminRepository(db)
	passwordUtil := password.NewPasswordUtil()
	tokenUtil := token.NewTokenUtil()
	cloudinaryInstance, _ := config.SetupCloudinary()
	cloudinaryService := cloudinary.NewCloudinaryService(cloudinaryInstance)

	adminUseCase := usecases.NewAdminUsecase(adminRepo, passwordUtil, cloudinaryService, tokenUtil)
	adminController := controllers.NewAdminController(adminUseCase, v)

	// Public routes
	g.POST("/admin/login", adminController.Login)
	// Protected routes
	g.Use(echojwt.WithConfig(token.GetJWTConfig()), middlewares.IsSuperAdmin)
	g.POST("/admin/register", adminController.Register)
	g.GET("/admin", adminController.GetAllAdmins)
	g.DELETE("/admin/:id", adminController.DeleteAdmin)
	g.PUT("/admin/:id", adminController.UpdateAdmin)
	g.GET("/admin/search", adminController.SearchAdminByUsername)
}
