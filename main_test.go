package main

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var app App

func TestMain(m *testing.M) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(".testing.env file could not be loaded", err)
	}

	app = NewApp()
	if err := app.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	); err != nil {
		log.Fatal("Failed to initialize the App", err)
	}

	RunMigrations(app)
	code := m.Run()
	RollbackMigrations(app)

	os.Exit(code)
}
