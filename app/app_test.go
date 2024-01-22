package app

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

var testApp App

func TestMain(m *testing.M) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(".testing.env file could not be loaded ", err)
	}

	testApp = NewApp()
	if err := testApp.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	); err != nil {
		log.Fatal("Failed to initialize the App ", err)
	}

	if err := RunMigrations(testApp); err != nil {
		log.Fatal("Failed to run migrations ", err)
	}

	code := m.Run()

	if err := RollbackMigrations(testApp); err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}

func clearTable(tableName string) error {
	if _, err := testApp.DB.Exec("DELETE FROM " + tableName); err != nil {
		return err
	}
	if _, err := testApp.DB.Exec("ALTER TABLE " + tableName + " AUTO_INCREMENT = 1"); err != nil {
		return err
	}
	return nil
}

func addShortKey(originalURL string) (string, error) {
	shortKey := generateShortKey()
	_, err := testApp.DB.Exec("INSERT INTO urls(orignal_url, short_key, expire_time) VALUES($1, $2, $3)", originalURL, shortKey, time.Now().AddDate(100, 0, 0))
	return shortKey, err
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	testApp.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
