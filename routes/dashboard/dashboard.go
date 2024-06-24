package dashboard

import (
	"kreasi-nusantara-api/controllers"
	"kreasi-nusantara-api/middlewares"
	"kreasi-nusantara-api/repositories"
	"kreasi-nusantara-api/usecases"
	"kreasi-nusantara-api/utils/token"
	"kreasi-nusantara-api/utils/validation"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitProductDashboard(g *echo.Group, db *gorm.DB, v *validation.Validator) {
	cartRepo := repositories.NewCartRepository(db)
	cartUseCase := usecases.NewCartUseCase(cartRepo)

	productDashboardRepository := repositories.NewProductDashboardRepository(db)
	productDashboardUseCase := usecases.NewProductDashboardUseCase(productDashboardRepository, cartUseCase)
	productDashboardController := controllers.NewProductDashboardController(productDashboardUseCase, v)

	g.Use(echojwt.WithConfig(token.GetJWTConfig()), middlewares.IsAdminOrSuperAdmin)

	g.GET("/products-report", productDashboardController.GetReportProducts)
	g.GET("/dashboard-header", productDashboardController.GetHeaderProduct)
	g.GET("/products-chart", productDashboardController.GetChartProduct)
	g.GET("/events-report", productDashboardController.GetEventReport)
	g.GET("/events-chart", productDashboardController.GetChartEvents)

}
