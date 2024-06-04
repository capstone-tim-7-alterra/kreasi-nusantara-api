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
	"kreasi-nusantara-api/utils/token"
	"kreasi-nusantara-api/utils/validation"

	echojwt "github.com/labstack/echo-jwt/v4"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitUserRoute(g *echo.Group, db *gorm.DB, v *validation.Validator) {
	userRepo := repositories.NewUserRepository(db)
	redisClient := redis.NewRedisClient()
	passwordUtil := password.NewPasswordUtil()
	otpUtil := otp.NewOTPUtil()
	emailUtil := email.NewEmailUtil()
	tokenUtil := token.NewTokenUtil()

	userUseCase := usecases.NewUserUseCase(userRepo, passwordUtil, *redisClient, otpUtil, emailUtil, tokenUtil)
	userController := controllers.NewUserController(userUseCase, v, tokenUtil)

	// Public routes
	g.POST("/register", userController.Register)
	g.POST("/verify-otp", userController.VerifyOTP)
	g.POST("/login", userController.Login)
	g.POST("/forgot-password", userController.ForgotPassword)
	g.POST("/reset-password", userController.ResetPassword)

	// Protected routes
	g.Use(echojwt.WithConfig(token.GetJWTConfig()))
	g.GET("/users/me", userController.GetProfile)
	g.PUT("/users/me", userController.UpdateProfile)
	g.DELETE("/users/me", userController.DeleteProfile)
}