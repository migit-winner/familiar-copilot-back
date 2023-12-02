package main

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var upgrader = &websocket.Upgrader{}

func handleWebSocket(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	fmt.Printf("websocket connected: %s\n", ws.RemoteAddr())

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Printf("websocket error: %s\n", err)
			return nil
		}

		fmt.Printf("websocket receive: %s\n", msg)

		err = ws.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Printf("websocket error: %s\n", err)
			return nil
		}
	}

	//return nil
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/ws", handleWebSocket)
	e.Logger.Fatal(e.Start(":8080"))
}
