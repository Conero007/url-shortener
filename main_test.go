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

	if err := makeTestDB(); err != nil {
		log.Fatal(err)
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

	if err := RunMigrations(app); err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	if err := RollbackMigrations(app); err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}

func makeTestDB() error {
	app.DB.Exec("CREATE DATABASE " + os.Getenv("DB_NAME"))

	if _, err := app.DB.Exec("USE " + os.Getenv("DB_NAME")); err != nil {
		return err
	}

	return nil
}
