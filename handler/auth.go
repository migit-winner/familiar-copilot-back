package handler

import (
	"familiar-copilot-back/domain"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")

	// TODO: データベースからユーザーを取得する
	fmt.Printf("Login: username=%s, password=%s\n", username, password)
	userID := 1 // 仮のユーザーID

	claims := &domain.JwtCustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(200, map[string]string{
		"token": t,
	})
}
