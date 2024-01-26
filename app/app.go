package app

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"sync"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

var App *AppConfig

type AppConfig struct {
	Router *mux.Router
	DB     *sql.DB
	Redis  *redis.Client

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

func (a *AppConfig) InitializeDB(addr, user, password, DBName string) error {
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

	return err
}

func (a *AppConfig) InitializeRoutes() {
	a.Router = mux.NewRouter()
	App.Router.HandleFunc("/shorten", HandleURLShortening).Methods(http.MethodPost)
	App.Router.HandleFunc("/{key}", HandleRedirectToOriginalURL).Methods(http.MethodGet)
}

func (a *AppConfig) InitializeRedis(addr, password string) error {
	a.Redis = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	_, err := a.Redis.Ping(context.Background()).Result()
	return err
}

func (a *AppConfig) Run(addr string) error {
	log.Printf("Starting Server at http://%s\n", addr)
	if err := http.ListenAndServe(addr, a.Router); err != nil {
		return err
	}
	return nil
}
