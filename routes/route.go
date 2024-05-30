package routes

import (
	"kreasi-nusantara-api/utils/validation"
	"kreasi-nusantara-api/routes/user"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitRoute(e *echo.Echo, db *gorm.DB, v *validation.Validator) {
	baseRoute := e.Group("/api/v1")

	userRoute := baseRoute.Group("")

	user.InitUserRoute(userRoute, db, v)
}
