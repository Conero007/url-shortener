package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var App *AppConfig

type AppConfig struct {
	Router *mux.Router
	DB     *sql.DB

	wg    *sync.WaitGroup
	debug bool
}

func NewApp(debug bool) *AppConfig {
	App = &AppConfig{
		debug: debug,
		wg:    &sync.WaitGroup{},
	}
	return App
}

func (a *AppConfig) Initialize(host, port, user, password, DBName string) error {
	addr := fmt.Sprintf("%s:%s", host, port)
	cfg := mysql.Config{
		User:      user,
		Passwd:    password,
		Net:       "tcp",
		Addr:      addr,
		ParseTime: true,
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
	InitializeRoutes()

	return err
}

func (a *AppConfig) Run(addr string) error {
	log.Printf("Starting Server at http://%s\n", addr)
	if err := http.ListenAndServe(addr, a.Router); err != nil {
		return err
	}
	return nil
}
