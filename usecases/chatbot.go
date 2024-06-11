package usecases

import (
	"kreasi-nusantara-api/drivers/openai"
	"kreasi-nusantara-api/drivers/redis"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	msg "kreasi-nusantara-api/constants/message"
)

type ChatBotUseCase interface {
	AnswerChat(c echo.Context, userId uuid.UUID, question string) (string, error)
}

type chatBotUseCase struct {
	openAIService openai.OpenAIClient
	redisClient   redis.RedisClient
}

func NewChatBotUseCase(openAIService openai.OpenAIClient, redisClient redis.RedisClient) *chatBotUseCase {
	return &chatBotUseCase{
		openAIService: openAIService,
		redisClient:   redisClient,
	}
}

func (cb *chatBotUseCase) AnswerChat(c echo.Context, userId uuid.UUID, question string) (string, error) {
	userIDstr := userId.String()
    key := "chat_history:" + userIDstr

    history, err := cb.redisClient.Get(key)
    if err != nil && err.Error() != "redis: nil" {
        return "", err
    }

    limitedHistory := limitHistory(history, 10)

    prompt := limitedHistory + "\nUser: " + question

    answer, err := cb.openAIService.AnswerChat(prompt)
    if err != nil {
        return msg.FAILED_ANSWER_CHAT, err
    }

    newHistory := updateHistory(limitedHistory, "User: "+question, "Assistant: "+answer)
    err = cb.redisClient.Set(key, newHistory, time.Hour*6)
    if err != nil {
        return msg.FAILED_ANSWER_CHAT, err
    }

    return answer, nil
}

func limitHistory(history string, limit int) string {
	lines := strings.Split(history, "\n")
	var limitedLines []string

	if len(lines) > limit*2 {
		limitedLines = lines[len(lines)-limit*2:]
	} else {
		limitedLines = lines
	}

	return strings.Join(limitedLines, "\n")
}

func updateHistory(history, userLine, assistantLine string) string {
	return history + "\n" + userLine + "\n" + assistantLine
}