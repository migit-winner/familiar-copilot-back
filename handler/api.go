package handler

import (
	"familiar-copilot-back/domain"
	"fmt"
	"net/http"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Hundler struct {
	db *sql.DB
}

func NewHundler() *Hundler {
	return &Hundler{}
}

func (h *Hundler) DBConnect() error {
	// DB接続
	dbconf := "user:password@tcp(db:3306)/FAMILIA_COPILOT?charset=utf8mb4"

	var err error
	h.db, err = sql.Open("mysql", dbconf)
	if err != nil {
		return err
	}

	err = h.db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (h *Hundler) CreateUaer(c echo.Context) error {
	// リクエストパラメータ取得
	var user User
	err := c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "パラメータが不正です")
	}

	err = h.db.Ping()
	if err != nil {
		return err
	}

	_, err = h.db.Exec("INSERT INTO users (name, password) VALUES (?, ?)", user.Name, user.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "ユーザー登録エラー :"+err.Error())
	}

	return c.JSON(http.StatusOK, "ユーザー登録完了")
}

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
