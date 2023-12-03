package handler

import (
	"familiar-copilot-back/domain"
	"familiar-copilot-back/infra"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleWebSocket(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &WebSocketHandler{
		upgrader: &websocket.Upgrader{},
		dbClient: &infra.DBClient{},
	}

	// JWTトークンを設定
	token := jwt.New(jwt.SigningMethodHS256)
	claims := &domain.JwtCustomClaims{
		UserID: 1,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 15000,
		},
	}
	token.Claims = claims
	c.Set("user", token)

	if assert.NoError(t, h.HandleWebSocket(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
