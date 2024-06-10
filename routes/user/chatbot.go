package user

import (
	"os"

	"github.com/labstack/echo/v4"

	"kreasi-nusantara-api/controllers"
	"kreasi-nusantara-api/usecases"
	oai "kreasi-nusantara-api/drivers/openai"
)

func InitChatBotRoute(g *echo.Group) {
	oaiKey := os.Getenv("OPENAI_API_KEY")

	oaiService := oai.NewOpenAIClient(oaiKey)

	chatBotUseCase := usecases.NewChatBotUseCase(oaiService)
	chatBotController := controllers.NewChatBotController(chatBotUseCase)

	g.GET("/chatbot", chatBotController.AnswerChat)
}