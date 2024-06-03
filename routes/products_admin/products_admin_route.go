package products_admin

import (
	"kreasi-nusantara-api/controllers"
	"kreasi-nusantara-api/repositories"
	"kreasi-nusantara-api/usecases"
	"kreasi-nusantara-api/utils/validation"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitProductAdminRoute(g *echo.Group, db *gorm.DB, v *validation.Validator) {
	productAdminRepository := repositories.NewProductAdminRepository(db)
	productAdminUsecase := usecases.NewProductAdminUseCase(productAdminRepository)
	productAdminController := controllers.NewProductsAdminController(productAdminUsecase, v)

	// g.GET("/products", productAdminController.GetAllProducts)
	// g.POST("/products", productAdminController.CreateProduct)
	// g.DELETE("/products/:id", productAdminController.DeleteProduct)
	// g.PUT("/products/:id", productAdminController.UpdateProduct)

	g.GET("/categories", productAdminController.GetAllCategories)
	g.POST("/categories", productAdminController.CreateCategory)
	g.DELETE("/categories/:id", productAdminController.DeleteCategory)
	g.PUT("/categories/:id", productAdminController.UpdateCategory)
}
