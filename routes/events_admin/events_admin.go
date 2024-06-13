package events_admin

// import (
// 	"kreasi-nusantara-api/config"
// 	"kreasi-nusantara-api/controllers"
// 	"kreasi-nusantara-api/drivers/cloudinary"
// 	"kreasi-nusantara-api/middlewares"
// 	"kreasi-nusantara-api/repositories"
// 	"kreasi-nusantara-api/usecases"
// 	"kreasi-nusantara-api/utils/token"
// 	"kreasi-nusantara-api/utils/validation"

// 	echojwt "github.com/labstack/echo-jwt/v4"
// 	"github.com/labstack/echo/v4"
// 	"gorm.io/gorm"
// )

// func InitEventsAdminRoute(g *echo.Group, db *gorm.DB, v *validation.Validator) {

// 	cloudinaryInstance, _ := config.SetupCloudinary()
// 	cloudinaryService := cloudinary.NewCloudinaryService(cloudinaryInstance)
// 	// tokenUtil := token.NewTokenUtil()

// 	eventAdminRepo := repositories.NewEventAdminRepository(db)
// 	eventAdminUsecase := usecases.NewEventAdminUseCase(eventAdminRepo)
// 	eventAdminController := controllers.NewEventsAdminController(eventAdminUsecase, v, cloudinaryService)

// 	g.Use(echojwt.WithConfig(token.GetJWTConfig()), middlewares.IsAdminOrSuperAdmin)
// 	g.GET("/events", eventAdminController.GetAllEvents)
// 	g.POST("/events", eventAdminController.CreateEventsAdmin)
// 	g.GET("/events/search", eventAdminController.SearchEventsAdmin)
// }
