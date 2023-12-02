package infra

import (
	"database/sql"
	"familiar-copilot-back/domain"

	_ "github.com/go-sql-driver/mysql"
)

type DBClient struct {
	db *sql.DB
}

func (c *DBClient) DBConnect() error {
	// DB接続
	dbconf := "user:password@tcp(db:3306)/FAMILIA_COPILOT?charset=utf8mb4"

	var err error
	c.db, err = sql.Open("mysql", dbconf)
	if err != nil {
		return err
	}

	err = c.db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (c *DBClient) CreateUaer(username, password string) error {
	err := c.db.Ping()
	if err != nil {
		return err
	}

	_, err = c.db.Exec("INSERT INTO users (name, password) VALUES (?, ?)", username, password)
	if err != nil {
		return err
	}

	return nil
}

func (c *DBClient) GetUserByID(userID int) (domain.User, error) {
	var user domain.User

	err := c.db.QueryRow("SELECT id, name, password FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Name, &user.Password)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (c *DBClient) GetUserByName(name string) (domain.User, error) {
	var user domain.User

	err := c.db.QueryRow("SELECT id, name, password FROM users WHERE name = ?", name).Scan(&user.ID, &user.Name, &user.Password)
	if err != nil {
		return user, err
	}

	return user, nil
}
