package main

import (
	"familiar-copilot-back/domain"
	"familiar-copilot-back/handler"

	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	hundler := handler.NewHundler()
	err := hundler.DBConnect()
	if err != nil {
		log.Fatal(err)
	}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(
		middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
			AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		}))
	e.POST("/user", hundler.CreateUaer)
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
