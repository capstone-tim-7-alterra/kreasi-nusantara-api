package controllers

import (
	"kreasi-nusantara-api/usecases"
	"kreasi-nusantara-api/utils/token"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type chatBotController struct {
	chatbotUseCase usecases.ChatBotUseCase
	upgrader       websocket.Upgrader
	token          token.TokenUtil
}

func NewChatBotController(chatbotUseCase usecases.ChatBotUseCase, token token.TokenUtil) *chatBotController {
	return &chatBotController{
		chatbotUseCase: chatbotUseCase,
		token: token,
	}
}

func (cbc *chatBotController) AnswerChat(c echo.Context) error {
	claims := cbc.token.GetClaims(c)

	ws, err := cbc.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
			break
		}

		res, err := cbc.chatbotUseCase.AnswerChat(c, claims.ID, string(msg))
		if err != nil {
			c.Logger().Error(err)
		}

		err = ws.WriteMessage(websocket.TextMessage, []byte(res))
		if err != nil {
			c.Logger().Error(err)
		}
	}

	return nil
}
