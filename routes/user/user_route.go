package user

import (
	"kreasi-nusantara-api/controllers"
	"kreasi-nusantara-api/drivers/redis"
	// "kreasi-nusantara-api/middlewares"
	"kreasi-nusantara-api/repositories"
	"kreasi-nusantara-api/usecases"
	"kreasi-nusantara-api/utils/email"
	"kreasi-nusantara-api/utils/otp"
	"kreasi-nusantara-api/utils/password"
	"kreasi-nusantara-api/utils/validation"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitUserRoute(g *echo.Group, db *gorm.DB, v *validation.Validator) {
	userRepo := repositories.NewUserRepository(db)
	redisClient := redis.NewRedisClient()
	passwordUtil := password.NewPasswordUtil()
	otpUtil := otp.NewOTPUtil()
	emailUtil := email.NewEmailUtil()

	userUseCase := usecases.NewUserUseCase(userRepo, passwordUtil, *redisClient, otpUtil, emailUtil)
	userController := controllers.NewUserController(userUseCase, v)

	// Public routes
	g.POST("/register", userController.Register)
	g.POST("/verify-otp", userController.VerifyOTP)

	// Protected routes
	// g.Use(middlewares.Auth)
	// g.GET("/profile", userController.GetProfile)
}
