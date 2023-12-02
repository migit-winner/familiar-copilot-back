package handler

import (
	"familiar-copilot-back/domain"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WebSocketHundler struct {
	upgrader *websocket.Upgrader
}

func NewWebSocketHundler() *WebSocketHundler {
	return &WebSocketHundler{&websocket.Upgrader{}}
}

func (w *WebSocketHundler) websocketLoop(ws *websocket.Conn) {
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

func (w *WebSocketHundler) HandleWebSocket(c echo.Context) error {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return echo.ErrUnauthorized
	}

	claims := user.Claims.(*domain.JwtCustomClaims)
	userID := claims.UserID
	fmt.Printf("user id: %d\n", userID)

	ws, err := w.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	fmt.Printf("websocket connected: %s\n", ws.RemoteAddr())

	ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Your user ID is %d", userID)))
	go w.websocketLoop(ws)

	return nil
}
