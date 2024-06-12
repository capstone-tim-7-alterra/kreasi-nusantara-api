package articles

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

func InitArticlesRoute(g *echo.Group, db *gorm.DB, v *validation.Validator) {
	articleRepo := repositories.NewArticleRepository(db)
	articleUseCase := usecases.NewArticleUseCase(articleRepo)

	tokenUtil := token.NewTokenUtil()

	articleController := controllers.NewArticleController(articleUseCase, v, tokenUtil)

	g.Use(echojwt.WithConfig(token.GetJWTConfig()))
	g.GET("/articles", articleController.GetArticles)
	g.GET("/articles/:article_id", articleController.GetArticleByID)
	g.GET("/articles/search", articleController.SearchArticles)
	g.GET("/articles/:article_id/comments", articleController.GetCommentsByArticleID)
	g.POST("/articles/:article_id/comments", articleController.AddCommentToArticle)
	g.POST("/articles/:article_id/comments/:comment_id/reply", articleController.ReplyToComment)
	g.POST("/articles/:article_id/like", articleController.LikeArticle)
	g.POST("/articles/:article_id/unlike", articleController.UnlikeArticle)
}
