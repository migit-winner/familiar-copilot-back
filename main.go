package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Hundler struct {
	upgrader *websocket.Upgrader
	db       *sql.DB
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

func (w *Hundler) websocketLoop(ws *websocket.Conn) {
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

func (w *Hundler) handleWebSocket(c echo.Context) error {
	ws, err := w.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	fmt.Printf("websocket connected: %s\n", ws.RemoteAddr())

	go w.websocketLoop(ws)

	return nil
}

func (h *Hundler) createUaer(c echo.Context) error {
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

func main() {
	e := echo.New()
	hundler := NewHundler()
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
	e.GET("/ws", hundler.handleWebSocket)
	e.POST("/user", hundler.createUaer)
	e.Logger.Fatal(e.Start(":8080"))
}
