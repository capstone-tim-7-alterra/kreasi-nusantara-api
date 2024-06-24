package articles_admin

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

func InitArticleAdminRoute(g *echo.Group, db *gorm.DB, v *validation.Validator) {

	cloudinaryInstance, _ := config.SetupCloudinary()
	cloudinaryService := cloudinary.NewCloudinaryService(cloudinaryInstance)
	tokenUtil := token.NewTokenUtil()
	adminRepo := repositories.NewAdminRepository(db)

	articleAdminRepo := repositories.NewArticleAdminRepository(db)
	articleAdminUsecase := usecases.NewArticleUseCaseAdmin(articleAdminRepo, tokenUtil, adminRepo)
	articleAdminController := controllers.NewArticlesAdminController(articleAdminUsecase, v, cloudinaryService, tokenUtil)

	g.Use(echojwt.WithConfig(token.GetJWTConfig()), middlewares.IsAdminOrSuperAdmin)
	g.GET("/articles", articleAdminController.GetArticles)
	g.POST("/articles", articleAdminController.CreateArticlesAdmin)
	g.DELETE("/articles/:id", articleAdminController.DeleteArticlesAdmin)
	g.PUT("/articles/:id", articleAdminController.UpdateArticlesAdmin)
	g.GET("/articles/search", articleAdminController.SearchArticles)
	g.GET("/articles/:id", articleAdminController.GetArticleByID)
}
