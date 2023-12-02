package handler

import (
	"familiar-copilot-back/domain"
	"familiar-copilot-back/infra"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WebSocketHandler struct {
	upgrader *websocket.Upgrader
	dbClient *infra.DBClient
}

func NewWebSocketHandler(dbClient *infra.DBClient) *WebSocketHandler {
	return &WebSocketHandler{&websocket.Upgrader{}, dbClient}
}

func (w *WebSocketHandler) websocketLoop(ws *websocket.Conn, user domain.User) {
	defer ws.Close()
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Printf("websocket error: %s\n", err)
			return
		}

		fmt.Printf("websocket receive: %s\n", msg)

		err = ws.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Printf("websocket error: %s\n", err)
			return
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
