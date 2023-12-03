package handler

import (
	"encoding/json"
	"familiar-copilot-back/domain"
	"familiar-copilot-back/infra"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WebSocketHandler struct {
	upgrader     *websocket.Upgrader
	dbClient     *infra.DBClient
	openAIClient *infra.OpenAIClient
}

func NewWebSocketHandler(dbClient *infra.DBClient, openAIClient *infra.OpenAIClient) *WebSocketHandler {
	return &WebSocketHandler{&websocket.Upgrader{}, dbClient, openAIClient}
}

type Message struct {
	Before string `json:"before"`
	After  string `json:"after"`
}

func (w *WebSocketHandler) websocketLoop(ws *websocket.Conn, user domain.User) {
	defer ws.Close()
	for {
		var message Message
		err := ws.ReadJSON(&message)

		if err == nil {
			middleText, err := w.openAIClient.GenMiddleText(message.Before, message.After)
			if err != nil {
				fmt.Printf("openai error: %s\n", err)
				var response struct {
					Error string `json:"error"`
				}
				response.Error = "openai error"
				responseJSON, _ := json.Marshal(response)

				err = ws.WriteMessage(websocket.TextMessage, responseJSON)
				if err != nil {
					fmt.Printf("websocket error: %s\n", err)
					return
				}
				continue
			}
			var response struct {
				Middle string `json:"middle"`
			}
			response.Middle = middleText
			responseJSON, _ := json.Marshal(response)
			err = ws.WriteMessage(websocket.TextMessage, responseJSON)
			if err != nil {
				fmt.Printf("websocket error: %s\n", err)
				return
			}
		} else {
			var response struct {
				Error string `json:"error"`
			}
			response.Error = "invalid json"
			responseJSON, _ := json.Marshal(response)
			err = ws.WriteMessage(websocket.TextMessage, responseJSON)
			if err != nil {
				fmt.Printf("websocket error: %s\n", err)
				return
			}
		}

	}
}

func (w *WebSocketHandler) HandleWebSocket(c echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return echo.ErrUnauthorized
	}

	claims := token.Claims.(*domain.JwtCustomClaims)
	userID := claims.UserID
	user, err := w.dbClient.GetUserByID(userID)
	if err != nil {
		return err
	}

	ws, err := w.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	fmt.Printf("websocket connected: %s\n", ws.RemoteAddr())

	ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Your user ID is %d", userID)))
	go w.websocketLoop(ws, user)

	return nil
}
