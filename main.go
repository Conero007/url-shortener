package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(".env file could not be loaded ", err)
	}

	app := NewApp()
	if err := app.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	); err != nil {
		log.Fatal("Failed to initialize the App ", err)
	}

	// addr := fmt.Sprintf("%s:%s", os.Getenv("ADDR"), os.Getenv("PORT"))
	// if err := app.Run(addr); err != nil {
	// 	log.Fatal("Failed to Run the APP ", err)
	// }
}
