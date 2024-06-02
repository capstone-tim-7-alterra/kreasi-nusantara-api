package routes

import (
	"kreasi-nusantara-api/utils/validation"
	"kreasi-nusantara-api/routes/user"
	"kreasi-nusantara-api/routes/admin"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitRoute(e *echo.Echo, db *gorm.DB, v *validation.Validator) {
	baseRoute := e.Group("/api/v1")

	userRoute := baseRoute.Group("")

	user.InitUserRoute(userRoute, db, v)
	admin.InitAdminRoute(baseRoute, db, v)
}
