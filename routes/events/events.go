package events

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

func InitEventsRoute(g *echo.Group, db *gorm.DB, v *validation.Validator) {
	eventRepo := repositories.NewEventRepository(db)
	eventUseCase := usecases.NewEventUseCase(eventRepo)
	eventController := controllers.NewEventController(eventUseCase, v)

	g.Use(echojwt.WithConfig(token.GetJWTConfig()))
	g.GET("/events", eventController.GetEvents)
	g.GET("/events/:event_id", eventController.GetEventByID)
	g.GET("/events/category/:category_id", eventController.GetEventsByCategory)
	g.GET("/events/search", eventController.SearchEvents)
}
