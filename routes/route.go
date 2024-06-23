package routes

import (
	"kreasi-nusantara-api/routes/admin"
	"kreasi-nusantara-api/routes/articles"
	"kreasi-nusantara-api/routes/articles_admin"
	"kreasi-nusantara-api/routes/cart"
	"kreasi-nusantara-api/routes/events"
	"kreasi-nusantara-api/routes/events_admin"
	"kreasi-nusantara-api/routes/event_transactions"
	"kreasi-nusantara-api/routes/product_transactions"
	"kreasi-nusantara-api/routes/products"
	"kreasi-nusantara-api/routes/products_admin"
	"kreasi-nusantara-api/routes/user"
	"kreasi-nusantara-api/routes/webhook"
	"kreasi-nusantara-api/utils/validation"
	"kreasi-nusantara-api/routes/product_dashboard"


	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitRoute(e *echo.Echo, db *gorm.DB, v *validation.Validator) {
	baseRoute := e.Group("/api/v1")

	userRoute := baseRoute.Group("")
	adminRoute := baseRoute.Group("")
	productsadminRoute := baseRoute.Group("/admin")
	productsRoute := baseRoute.Group("")
	eventsRoute := baseRoute.Group("")
	chatBotRoute := baseRoute.Group("")
	eventsadminRoute := baseRoute.Group("/admin")
	articlesAdminRoute := baseRoute.Group("/admin")
	cartRoute := baseRoute.Group("")
	productTransactionRoute := baseRoute.Group("")
	eventTransactionRoute := baseRoute.Group("")
	paymentNotifRoute := baseRoute.Group("")
	productDashboardRoute := baseRoute.Group("/admin")

	user.InitUserRoute(userRoute, db, v)
	user.InitUserAddressesRoute(userRoute, db, v)
	user.InitChatBotRoute(chatBotRoute)
	admin.InitAdminRoute(adminRoute, db, v)
	products_admin.InitProductAdminRoute(productsadminRoute, db, v)
	products.InitProductsRoute(productsRoute, db, v)
	articles.InitArticlesRoute(baseRoute, db, v)
	articles_admin.InitArticleAdminRoute(articlesAdminRoute, db, v)
	events.InitEventsRoute(eventsRoute, db, v)
	events_admin.InitEventsAdminRoute(eventsadminRoute, db, v)
	cart.InitCartRoute(cartRoute, db, v)
	product_transactions.InitProductTransactionsRoute(productTransactionRoute, db, v)
	event_transactions.InitEventTransactionsRoute(eventTransactionRoute, db, v)
	webhook.InitWebhookRoute(paymentNotifRoute, db)
	product_dashboard.InitProductDashboard(productDashboardRoute, db, v)
}
