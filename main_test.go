package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/Conero007/url-shortener-golang/app"
	"github.com/Conero007/url-shortener-golang/constants"
	"github.com/Conero007/url-shortener-golang/database"
	"github.com/joho/godotenv"
)

var TestApp *app.AppConfig

func TestMain(m *testing.M) {
	if err := godotenv.Load(".testing.env"); err != nil {
		log.Fatal(".testing.env file could not be loaded ", err)
	}

	TestApp = app.NewApp(true)

	if err := TestApp.InitializeDB(
		os.Getenv("DB_ADDR"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	); err != nil {
		log.Fatal("Failed to initialize db ", err)
	}

	TestApp.InitializeRoutes()

	if err := TestApp.InitializeRedis(
		os.Getenv("REDIS_ADDR"),
		os.Getenv("REDIS_PASSWORD"),
	); err != nil {
		log.Fatal("Failed to initialize redis ", err)
	}

	m.Run()

	if err := database.RollbackMigrations(TestApp.DB, "create_urls_table"); err != nil {
		log.Fatal(err)
	}
}

func TestCreateShortenURL(t *testing.T) {
	if err := clearData("urls"); err != nil {
		t.Errorf("Could not clear urls table. ERROR: %s", err.Error())
		return
	}

	var jsonStr = []byte(`{"url":"https://www.google.com/"}`)
	req, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if !validateShortenAPIResponse(t, m) {
		return
	}
}

func TestCreateShortenURLWithEmptyPayload(t *testing.T) {
	if err := clearData("urls"); err != nil {
		t.Errorf("Could not clear urls table. ERROR: %s", err.Error())
		return
	}

	response := sendRequesttoShortenAPI(`{}`)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "URL not given" {
		t.Errorf("Expected the 'error' key of the response to be set to 'URL not given'. Got '%s'", m["error"])
	}
}

func TestCreateShortenURLWithInvalidURL(t *testing.T) {
	if err := clearData("urls"); err != nil {
		t.Errorf("Could not clear urls table. ERROR: %s", err.Error())
		return
	}

	response := sendRequesttoShortenAPI(`{"url":"google.com"}`)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Invalid URL given" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Invalid URL given'. Got '%s'", m["error"])
	}
}

func TestCreateShortenURLWithCustomURL(t *testing.T) {
	if err := clearData("urls"); err != nil {
		t.Errorf("Could not clear urls table. ERROR: %s", err.Error())
		return
	}

	response := sendRequesttoShortenAPI(`{"url":"https://www.google.com/", "custom_short_key": "654321"}`)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if !validateShortenAPIResponse(t, m) {
		return
	}

	shortURL := m["short_url"].(string)
	shortKey := shortURL[len(shortURL)-constants.SHORT_KEY_LENGTH:]

	if shortKey != "654321" {
		t.Errorf("Expected short key to be (654321), got %s", shortKey)
	}
}

func TestCreateShortenURLWithInvalidCustomURL_SpecialCharValidation(t *testing.T) {
	if err := clearData("urls"); err != nil {
		t.Errorf("Could not clear urls table. ERROR: %s", err.Error())
		return
	}

	response := sendRequesttoShortenAPI(`{"url":"https://www.google.com/", "custom_short_key": "65#321"}`)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Invalid custom short key" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Invalid custom short key'. Got '%s'", m["error"])
	}
}

func TestCreateShortenURLWithInvalidCustomURL_LengthValidation(t *testing.T) {
	if err := clearData("urls"); err != nil {
		t.Errorf("Could not clear urls table. ERROR: %s", err.Error())
		return
	}

	response := sendRequesttoShortenAPI(`{"url":"https://www.google.com/", "custom_short_key": "1234567"}`)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Invalid custom short key" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Invalid custom short key'. Got '%s'", m["error"])
	}
}

func TestRedirectViaShortKey(t *testing.T) {
	if err := clearData("urls"); err != nil {
		t.Errorf("Could not clear urls table. ERROR: %s", err.Error())
		return
	}

	response := sendRequesttoShortenAPI(`{"url": "https://www.google.com/"}`)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if !validateShortenAPIResponse(t, m) {
		return
	}

	shortURL := m["short_url"].(string)
	shortKey := shortURL[len(shortURL)-constants.SHORT_KEY_LENGTH:]

	req, _ := http.NewRequest("GET", "/"+shortKey, nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusMovedPermanently, response.Code)

	if response.Result().Header.Get("Location") != "https://www.google.com/" {
		t.Error("Excepted redirect url https://www.google.com/, found ", response.Result().Header.Get("Location"))
	}
}

func TestGetNonExistentShortKey(t *testing.T) {
	if err := clearData("urls"); err != nil {
		t.Errorf("Could not clear urls table. ERROR: %s", err.Error())
		return
	}

	req, _ := http.NewRequest("GET", "/123456", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Short Key not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Short Key not found'. Got '%s'", m["error"])
	}
}

func TestGetExpiredShortKey(t *testing.T) {
	if err := clearData("urls"); err != nil {
		t.Errorf("Could not clear urls table. ERROR: %s", err.Error())
		return
	}

	shortKey, err := addShortKey("https://www.google.com/", time.Now().Add(-time.Second))
	if err != nil {
		t.Errorf("Failed to add original url to DB. ERROR: %s", err.Error())
		return
	}

	req, _ := http.NewRequest("GET", "/"+shortKey, nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Short Key not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Short Key not found'. Got '%s'", m["error"])
	}

	originalURL := fetchOriginalURL(shortKey)
	if originalURL != "" {
		t.Error("Expected the expired short key to be deleted from the DB, but found it")
	}
}

func TestInvalidShortKey_LengthValidation(t *testing.T) {
	if err := clearData("urls"); err != nil {
		t.Errorf("Could not clear urls table. ERROR: %s", err.Error())
		return
	}

	req, _ := http.NewRequest("GET", "/1234567", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Invalid short key" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Invalid short key'. Got '%s'", m["error"])
	}
}

func TestInvalidShortKey_SpecialCharValidation(t *testing.T) {
	if err := clearData("urls"); err != nil {
		t.Errorf("Could not clear urls table. ERROR: %s", err.Error())
		return
	}

	req, _ := http.NewRequest("GET", "/123#56", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Invalid short key" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Invalid short key'. Got '%s'", m["error"])
	}
}

func TestDuplicateOriginalURL(t *testing.T) {
	if err := clearData("urls"); err != nil {
		t.Errorf("Could not clear urls table. ERROR: %s", err.Error())
		return
	}

	response1 := sendRequesttoShortenAPI(`{"url":"https://www.google.com/"}`)
	checkResponseCode(t, http.StatusCreated, response1.Code)

	var m1 map[string]interface{}
	json.Unmarshal(response1.Body.Bytes(), &m1)
	if !validateShortenAPIResponse(t, m1) {
		return
	}

	response2 := sendRequesttoShortenAPI(`{"url":"https://www.google.com/"}`)
	checkResponseCode(t, http.StatusCreated, response2.Code)

	var m2 map[string]interface{}
	json.Unmarshal(response2.Body.Bytes(), &m2)
	if !validateShortenAPIResponse(t, m1) {
		return
	}

	if m1["short_url"].(string) == m2["short_url"].(string) {
		t.Errorf("Expected different short url for duplicate request for the same origianl url")
	}
}

func clearData(tableName string) error {
	if _, err := TestApp.DB.Exec("DELETE FROM " + tableName); err != nil {
		return err
	}
	if _, err := TestApp.DB.Exec("ALTER TABLE " + tableName + " AUTO_INCREMENT = 1"); err != nil {
		return err
	}
	if _, err := TestApp.Redis.FlushAllAsync(context.Background()).Result(); err != nil {
		return err
	}
	return nil
}

func addShortKey(originalURL string, expireTime time.Time) (string, error) {
	shortKey := "123456"
	_, err := TestApp.DB.Exec("INSERT INTO urls(original_url, short_key, expire_time) VALUES(?, ?, ?)", originalURL, shortKey, expireTime)
	return shortKey, err
}

func fetchOriginalURL(shortKey string) string {
	var originalURL string
	TestApp.DB.QueryRow("SELECT original_url FROM urls WHERE short_key = ? LIMIT 1", shortKey).Scan(&originalURL)
	return originalURL
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	TestApp.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func validateShortenAPIResponse(t *testing.T, m map[string]interface{}) bool {
	if _, ok := m["original_url"]; !ok {
		t.Error("original_url field missing in the response")
		return false
	} else if m["original_url"] != "https://www.google.com/" {
		t.Error("original_url different from the one in the request")
		return false
	}

	if _, ok := m["short_url"]; !ok {
		t.Error("short_url field missing in the response")
		return false
	}

	if _, ok := m["expire_time"]; !ok {
		t.Error("expire_time field missing in the response")
		return false
	}

	if _, ok := m["short_url"].(string); !ok {
		t.Error("Failed to typecast short_url field to string.")
		return false
	}

	regexPattern := fmt.Sprintf(`^http://%s:%s/[A-Z a-z 0-9]{%d}$`, os.Getenv("APP_URL"), os.Getenv("PORT"), constants.SHORT_KEY_LENGTH)
	if ok, _ := regexp.MatchString(regexPattern, m["short_url"].(string)); !ok {
		t.Errorf("Expected short_url format to be 'http://%s:%s/xxxxxx'. Got '%v'", os.Getenv("APP_URL"), os.Getenv("PORT"), m["short_url"])
		return false
	}

	return true
}

func sendRequesttoShortenAPI(paylaod string) *httptest.ResponseRecorder {
	var jsonStr1 = []byte(paylaod)
	req1, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonStr1))
	req1.Header.Set("Content-Type", "application/json")
	response := executeRequest(req1)
	return response
}
