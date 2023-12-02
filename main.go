package main

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type WebSocketHundler struct {
	upgrader *websocket.Upgrader
}

func NewWebSocketHundler() *WebSocketHundler {
	return &WebSocketHundler{}
}

func websocketLoop(ws *websocket.Conn) {
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

func (h *WebSocketHundler) handleWebSocket(c echo.Context) error {
	ws, err := h.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	fmt.Printf("websocket connected: %s\n", ws.RemoteAddr())

	go websocketLoop(ws)

	return nil
}

func main() {
	e := echo.New()
	webSocketHundler := NewWebSocketHundler()
	e.Use(middleware.Logger())
	e.GET("/ws", webSocketHundler.handleWebSocket)
	e.Logger.Fatal(e.Start(":8080"))
}
