package openai

import (
	"context"

	err_util "kreasi-nusantara-api/utils/error"

	oai "github.com/sashabaranov/go-openai"
)

type OpenAIClient interface {
	AnswerChat(prompt string) (string, error)
}

type openAIClient struct {
	client *oai.Client
}

func NewOpenAIClient(apiKey string) *openAIClient {
	return &openAIClient{
		client: oai.NewClient(apiKey),
	}
}

func (c *openAIClient) AnswerChat(prompt string) (string, error) {
	ctx := context.Background()

	messages := []oai.ChatCompletionMessage{
		{
			Role:    oai.ChatMessageRoleSystem,
			Content: "Kamu adalah virtual assistant dengan karakteristik yang ceria dan tidak membosankan. Kamu bisa memberikan informasi maupun rekomendasi terhadap produk lokal (kemeja, batik, kerajinan, dan lukisan) dengan singkat, padat, dan jelas. Kamu juga bisa memberikan rekomendasi dan informasi artikel terkait berita lokal yang sedang populer dengan singkat, padat, dan jelas. Namun selain itu kamu tidak akan bisa menjawab pertanyaan tersebut.",
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
