package user

import (
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	"kreasi-nusantara-api/controllers"
	oai "kreasi-nusantara-api/drivers/openai"
	"kreasi-nusantara-api/drivers/redis"
	"kreasi-nusantara-api/usecases"
	"kreasi-nusantara-api/utils/token"
)

func InitChatBotRoute(g *echo.Group) {
	oaiKey := os.Getenv("OPENAI_API_KEY")

	oaiService := oai.NewOpenAIClient(oaiKey)
	redisClient := redis.NewRedisClient()
	tokenUtil := token.NewTokenUtil()

	chatBotUseCase := usecases.NewChatBotUseCase(oaiService, *redisClient)
	chatBotController := controllers.NewChatBotController(chatBotUseCase, tokenUtil)

	g.Use(echojwt.WithConfig(token.GetJWTConfig()))
	g.GET("/chatbot", chatBotController.AnswerChat)
}
