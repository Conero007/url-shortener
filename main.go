package main

import (
	"log"
	"os"

	"github.com/Conero007/url-shortener-golang/app"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(".env file could not be loaded ", err)
	}

	app := app.NewApp(false)
	if err := app.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	); err != nil {
		log.Fatal("Failed to initialize the App ", err)
	}

	if err := app.Run(":" + os.Getenv("PORT")); err != nil {
		log.Fatal("Failed to Run the APP ", err)
	}
}
