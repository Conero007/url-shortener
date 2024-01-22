package app

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func NewApp() App {
	return App{}
}

func (a *App) Initialize(host, port, user, password, DBName string) error {
	addr := fmt.Sprintf("%s:%s", host, port)
	cfg := mysql.Config{
		User:   user,
		Passwd: password,
		Net:    "tcp",
		Addr:   addr,
		DBName: DBName,
	}

	var err error
	a.DB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return err
	}

	pingErr := a.DB.Ping()
	if pingErr != nil {
		return err
	}

	a.Router = mux.NewRouter()

	return err
}

func (a *App) Run(addr string) error {
	return nil
}
