package handler

import (
	"familiar-copilot-back/domain"
	"familiar-copilot-back/infra"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockDBClient struct{}

func (m *MockDBClient) CreateUser(name string, password string) error {
	return nil
}

func (m *MockDBClient) GetUserByName(name string) (domain.User, error) {
	return domain.User{Name: "test", Password: "test"}, nil
}

func TestAPIHandler_CreateUser(t *testing.T) {
	h := &APIHandler{dbClient: &MockDBClient{}}
	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"test","password":"test"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, h.CreateUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `"ユーザー登録完了"`, rec.Body.String())
	}
}

func TestAPIHandler_Login(t *testing.T) {
	h := &APIHandler{dbClient: &MockDBClient{}}
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/?username=test&password=test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, h.Login(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "token")
	}
}
