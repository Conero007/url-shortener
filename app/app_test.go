package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.testing.env"); err != nil {
		log.Fatal(".testing.env file could not be loaded ", err)
	}

	testApp := NewApp(true)
	if err := testApp.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	); err != nil {
		log.Fatal("Failed to initialize the App ", err)
	}

	code := m.Run()

	if err := RollbackMigrations(); err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}

func TestCreateShortenURL(t *testing.T) {
	if err := clearTable("urls"); err != nil {
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

	if _, ok := m["original_url"]; !ok {
		t.Error("original_url field missing in the response")
		return
	} else if m["original_url"] != "https://www.google.com/" {
		t.Error("original_url different from the one in the request")
	}

	if _, ok := m["short_url"]; !ok {
		t.Error("short_url field missing in the response")
		return
	}

	if _, ok := m["expire_time"]; !ok {
		t.Error("expire_time field missing in the response")
		return
	}

	if _, ok := m["short_url"].(string); !ok {
		t.Error("Failed to typecast short_url field to string.")
		return
	}

	regexPattern := fmt.Sprintf(`^http://%s:%s/[A-Z a-z 0-9]{6}$`, os.Getenv("APP_URL"), os.Getenv("PORT"))
	if ok, _ := regexp.MatchString(regexPattern, m["short_url"].(string)); !ok {
		t.Errorf("Expected short_url format to be 'http://%s:%s/xxxxxx'. Got '%v'", os.Getenv("APP_URL"), os.Getenv("PORT"), m["short_url"])
	}
}

func TestShortenURLNotGivenValidation(t *testing.T) {
	if err := clearTable("urls"); err != nil {
		t.Errorf("Could not clear urls table. ERROR: %s", err.Error())
		return
	}

	var jsonStr = []byte(`{}`)
	req, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "URL not given" {
		t.Errorf("Expected the 'error' key of the response to be set to 'URL not given'. Got '%s'", m["error"])
	}
}

func TestShortenInvalidURLValidation(t *testing.T) {
	if err := clearTable("urls"); err != nil {
		t.Errorf("Could not clear urls table. ERROR: %s", err.Error())
		return
	}

	var jsonStr = []byte(`{"url":"asdasdffsadklj"}`)
	req, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Invalid URL given" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Invalid URL given'. Got '%s'", m["error"])
	}
}

func TestRedirectViaShortKey(t *testing.T) {
	if err := clearTable("urls"); err != nil {
		t.Errorf("Could not clear urls table. ERROR: %s", err.Error())
		return
	}

	shortKey, err := addShortKey("https://www.google.com/", time.Now().Add(60*time.Second))
	if err != nil {
		t.Errorf("Failed to add original url to DB. ERROR: %s", err.Error())
		return
	}

	req, _ := http.NewRequest("GET", "/"+shortKey, nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusMovedPermanently, response.Code)

	if response.Result().Header.Get("Location") != "https://www.google.com/" {
		t.Error("Excepted redirect url https://www.google.com/, found ", response.Result().Header.Get("Location"))
	}
}

func TestGetNonExistentShortKey(t *testing.T) {
	if err := clearTable("urls"); err != nil {
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
	if err := clearTable("urls"); err != nil {
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

func TestShortKeyLengthVaidation(t *testing.T) {
	if err := clearTable("urls"); err != nil {
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

func TestShortKeySpecialCharVaidation(t *testing.T) {
	if err := clearTable("urls"); err != nil {
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

func clearTable(tableName string) error {
	if _, err := App.DB.Exec("DELETE FROM " + tableName); err != nil {
		return err
	}
	if _, err := App.DB.Exec("ALTER TABLE " + tableName + " AUTO_INCREMENT = 1"); err != nil {
		return err
	}
	return nil
}

func addShortKey(originalURL string, expireTime time.Time) (string, error) {
	shortKey := generateShortKey()
	_, err := App.DB.Exec("INSERT INTO urls(original_url, short_key, expire_time) VALUES(?, ?, ?)", originalURL, shortKey, expireTime)
	return shortKey, err
}

func fetchOriginalURL(shortKey string) string {
	var originalURL string
	App.DB.QueryRow("SELECT original_url FROM urls WHERE short_key = ? LIMIT 1", shortKey).Scan(&originalURL)
	return originalURL
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	App.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
