package openai

import (
	"context"

	err_util "kreasi-nusantara-api/utils/error"

	oai "github.com/sashabaranov/go-openai"
)

type OpenAIClient interface {
	AnswerChat(prompt, instruction string) (string, error)
}

type openAIClient struct {
	client *oai.Client
}

func NewOpenAIClient(apiKey string) *openAIClient {
	return &openAIClient{
		client: oai.NewClient(apiKey),
	}
}

func (c *openAIClient) AnswerChat(prompt, instruction string) (string, error) {
	ctx := context.Background()

	messages := []oai.ChatCompletionMessage{
		{
			Role:    oai.ChatMessageRoleSystem,
			Content: instruction,
		},
		{
			Role:    oai.ChatMessageRoleUser,
			Content: prompt,
		},
	}

	req := oai.ChatCompletionRequest{
		Model:    oai.GPT3Dot5Turbo,
		Messages: messages,
	}

	response, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err_util.ErrExternalServiceError
	}

	return response.Choices[0].Message.Content, nil
}
