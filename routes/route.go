package routes

import (
	"kreasi-nusantara-api/routes/admin"
	"kreasi-nusantara-api/routes/articles"
	"kreasi-nusantara-api/routes/articles_admin"
	"kreasi-nusantara-api/routes/events"
	// "kreasi-nusantara-api/routes/events_admin"
	"kreasi-nusantara-api/routes/products"
	"kreasi-nusantara-api/routes/products_admin"
	"kreasi-nusantara-api/routes/user"
	"kreasi-nusantara-api/utils/validation"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitRoute(e *echo.Echo, db *gorm.DB, v *validation.Validator) {
	baseRoute := e.Group("/api/v1")
	baseAdminRoute := e.Group("/api/v1/admin")

	userRoute := baseRoute.Group("")
	adminRoute := baseRoute.Group("")
	productsadminRoute := baseRoute.Group("")
	productsRoute := baseRoute.Group("")
	eventsRoute := baseRoute.Group("")
	chatBotRoute := baseRoute.Group("")
	// eventsadminRoute := baseRoute.Group("")
	articlesAdminRoute := baseAdminRoute.Group("")

	user.InitUserRoute(userRoute, db, v)
	user.InitUserAddressesRoute(userRoute, db, v)
	user.InitChatBotRoute(chatBotRoute)
	admin.InitAdminRoute(adminRoute, db, v)
	products_admin.InitProductAdminRoute(productsadminRoute, db, v)
	products.InitProductsRoute(productsRoute, db, v)
	articles.InitArticlesRoute(baseRoute, db, v)
	articles_admin.InitArticleAdminRoute(articlesAdminRoute, db, v)
	events.InitEventsRoute(eventsRoute, db, v)
}
	// events_admin.InitEventsAdminRoute(eventsadminRoute, db, v)
