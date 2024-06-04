package products_admin

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

func InitProductAdminRoute(g *echo.Group, db *gorm.DB, v *validation.Validator) {

	cloudinaryInstance, _ := config.SetupCloudinary()
	cloudinaryService := cloudinary.NewCloudinaryService(cloudinaryInstance)
	tokenUtil := token.NewTokenUtil()

	productAdminRepository := repositories.NewProductAdminRepository(db)
	productAdminUsecase := usecases.NewProductAdminUseCase(productAdminRepository, cloudinaryService, tokenUtil)
	productAdminController := controllers.NewProductsAdminController(productAdminUsecase, v)


	
	// g.DELETE("/products/:id", productAdminController.DeleteProduct)
	// g.PUT("/products/:id", productAdminController.UpdateProduct)

	g.Use(echojwt.WithConfig(token.GetJWTConfig()), middlewares.IsAdminOrSuperAdmin)
	g.GET("/categories", productAdminController.GetAllCategories)
	g.POST("/categories", productAdminController.CreateCategory)
	g.DELETE("/categories/:id", productAdminController.DeleteCategory)
	g.PUT("/categories/:id", productAdminController.UpdateCategory)
	g.POST("/products", productAdminController.CreateProduct)
	g.GET("/products", productAdminController.GetAllProducts)
	g.DELETE("/products/:id", productAdminController.DeleteProduct)
	g.PUT("/products/:id", productAdminController.UpdateProduct)
	g.GET("/products/search", productAdminController.SearchProductByName)
}
