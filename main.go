package main

import (
	"familiar-copilot-back/domain"
	"familiar-copilot-back/handler"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/login", handler.Login)

	restricted := e.Group("")
	restricted.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte("secret"),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return &domain.JwtCustomClaims{}
		},
	}))
	webSocketHundler := handler.NewWebSocketHundler()
	restricted.GET("/ws", webSocketHundler.HandleWebSocket)

	e.Logger.Fatal(e.Start(":8080"))
}
