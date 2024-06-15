package events_admin

import (
	"kreasi-nusantara-api/config"
	"kreasi-nusantara-api/controllers"
	"kreasi-nusantara-api/drivers/cloudinary"
	"kreasi-nusantara-api/middlewares"
	"kreasi-nusantara-api/repositories"
	"kreasi-nusantara-api/usecases"
	"kreasi-nusantara-api/utils/token"
	"kreasi-nusantara-api/utils/validation"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitEventsAdminRoute(g *echo.Group, db *gorm.DB, v *validation.Validator) {

	apiKey := "0ffe6a2a51f6dcee7cfa6b34f4649ee03ed9ec8a1cf062291abfb7aaeaf5067e"
	
	cloudinaryInstance, _ := config.SetupCloudinary()
	cloudinaryService := cloudinary.NewCloudinaryService(cloudinaryInstance)
	// tokenUtil := token.NewTokenUtil()
	wilayahUsecase := usecases.NewRegionUseCase(apiKey)
	wilayahController := controllers.NewRegionController(wilayahUsecase)

	eventAdminRepo := repositories.NewEventAdminRepository(db)
	eventAdminUsecase := usecases.NewEventAdminUseCase(eventAdminRepo)
	eventAdminController := controllers.NewEventsAdminController(eventAdminUsecase, v, cloudinaryService)

	g.Use(echojwt.WithConfig(token.GetJWTConfig()), middlewares.IsAdminOrSuperAdmin)
	g.GET("/events", eventAdminController.GetAllEvents)
	g.POST("/events", eventAdminController.CreateEventsAdmin)
	g.GET("/events/search", eventAdminController.SearchEventsAdmin)
	g.GET("/events/:event_id", eventAdminController.GetEventByID)
	g.PUT("/events/:event_id", eventAdminController.UpdateEventsAdmin)
	g.DELETE("/events/:event_id", eventAdminController.DeleteEventsAdmin)

	g.POST("/events/categories", eventAdminController.CreateCategoriesEvent)
	g.GET("/events/categories", eventAdminController.GetCategoriesEvent)
	g.PUT("/events/categories/:id", eventAdminController.UpdateCategoriesEvent)
	g.DELETE("/events/categories/:id", eventAdminController.DeleteCategoriesEvent)

	g.GET("/events/provinces", wilayahController.GetProvincesHandler)
	g.GET("/events/districts", wilayahController.GetDistrictsHandler)
	g.GET("/events/subdistricts", wilayahController.GetSubdistrictsHandler)

	g.GET("/events/ticket-types", eventAdminController.GetTicketType)
	g.POST("/events/ticket-types", eventAdminController.CreateTicketType)
	g.DELETE("/events/ticket-types/:id", eventAdminController.DeleteTicketType)

	g.GET("/events/:event_id/prices", eventAdminController.GetPricesByEventID)
	g.GET("/prices/:price_id", eventAdminController.GetDetailPrices)
	g.DELETE("/prices/:price_id", eventAdminController.DeletePrices)
	g.PUT("/prices/:price_id", eventAdminController.UpdatePrices)


}
