package controllers

import (
	"kreasi-nusantara-api/usecases"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type chatBotController struct {
	chatbotUseCase usecases.ChatBotUseCase
	upgrader       websocket.Upgrader
}

func NewChatBotController(chatbotUseCase usecases.ChatBotUseCase) *chatBotController {
	return &chatBotController{
		chatbotUseCase: chatbotUseCase,
	}
}

func (cbc *chatBotController) AnswerChat(c echo.Context) error {
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

		res, err := cbc.chatbotUseCase.AnswerChat(c, string(msg))
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