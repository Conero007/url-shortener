package app

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var App *AppConfig

type AppConfig struct {
	Router *mux.Router
	DB     *sql.DB
}

func NewApp() *AppConfig {
	App = &AppConfig{}
	return App
}

func (a *AppConfig) Initialize(host, port, user, password, DBName string) error {
	addr := fmt.Sprintf("%s:%s", host, port)
	cfg := mysql.Config{
		User:   user,
		Passwd: password,
		Net:    "tcp",
		Addr:   addr,
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

	if err := RunMigrations(DBName); err != nil {
		return err
	}

	a.Router = mux.NewRouter()

	return err
}

func (a *AppConfig) Run(addr string) error {
	return nil
}