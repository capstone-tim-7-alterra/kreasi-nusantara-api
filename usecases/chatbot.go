package usecases

import (
	"kreasi-nusantara-api/drivers/openai"

	"github.com/labstack/echo/v4"
)

type ChatBotUseCase interface {
	AnswerChat(c echo.Context, question string) (string, error)
}

type chatBotUseCase struct {
	openAIService openai.OpenAIClient
}

func NewChatBotUseCase(openAIService openai.OpenAIClient) *chatBotUseCase {
	return &chatBotUseCase{
		openAIService: openAIService,
	}
}

func (cb *chatBotUseCase) AnswerChat(c echo.Context, question string) (string, error) {
	answer, err := cb.openAIService.AnswerChat(question)
	if err != nil {
		return "Maaf saya bodoh", err
	}

	return answer, nil
}