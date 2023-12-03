package handler

import (
	"familiar-copilot-back/domain"
	"familiar-copilot-back/infra"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type APIHandler struct {
	dbClient      *infra.DBClient
	openApiClient *infra.OpenAIClient
}

func NewAPIHandler(dbClient *infra.DBClient, openApiClient *infra.OpenAIClient) *APIHandler {
	return &APIHandler{dbClient, openApiClient}
}

func (h *APIHandler) CreateUaer(c echo.Context) error {
	// リクエストパラメータ取得
	var user domain.User
	err := c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "パラメータが不正です")
	}

	err = h.dbClient.CreateUaer(user.Name, user.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "ユーザー登録エラー :"+err.Error())
	}

	return c.JSON(http.StatusOK, "ユーザー登録完了")
}

func (h *APIHandler) Login(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")

	// TODO: データベースからユーザーを取得する
	user, err := h.dbClient.GetUserByName(username)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "ユーザーが存在しません")
	}
	if user.Password != password {
		return c.JSON(http.StatusUnauthorized, "パスワードが間違っています")
	}

	claims := &domain.JwtCustomClaims{
		UserID: user.ID,
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

func (h *APIHandler) GenMiddleText(c echo.Context) error {
	var reqBody struct {
		Before string `json:"before"`
		After  string `json:"after"`
	}

	err := c.Bind(&reqBody)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "パラメータが不正です")
	}

	middleText, err := h.openApiClient.GenMiddleText(reqBody.Before, reqBody.After)

	return c.JSON(200, map[string]string{
		"middle": middleText,
	})
}
