package app

import (
	"bytes"
	"encoding/json"
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

	testApp := NewApp()
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

func TestGetNonExistentShortKey(t *testing.T) {
	if err := clearTable("urls"); err != nil {
		t.Errorf("Could not clear urls table. ERROR: %s", err.Error())
		return
	}

	req, _ := http.NewRequest("GET", "/shorten/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Short Key not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Short Key not found'. Got '%s'", m["error"])
	}
}

func TestCreateShortKey(t *testing.T) {
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

	if _, ok := m["short_key"]; !ok {
		t.Error("short_key field missing in the response")
		return
	}

	if _, ok := m["expire_time"]; !ok {
		t.Error("expire_time field missing in the response")
		return
	}

	if _, ok := m["shorten_url"].(string); !ok {
		t.Error("Failed to typecast shorten_url field to string.")
		return
	}

	regexPattern := `^https://localhost:8080/[A-Z a-z 0-9]{6}$`
	if ok, _ := regexp.MatchString(regexPattern, m["shorten_url"].(string)); !ok {
		t.Errorf("Expected product name to be 'test product'. Got '%v'", m["name"])
	}
}

func TestRedirectViaShortKey(t *testing.T) {
	clearTable("url")

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

func TestExpiredShortKey(t *testing.T) {
	clearTable("urls")
	shortKey, err := addShortKey("https://www.google.com/", time.Now())
	if err != nil {
		t.Errorf("Failed to add original url to DB. ERROR: %s", err.Error())
		return
	}

	time.Sleep(time.Second)

	req, _ := http.NewRequest("GET", "/"+shortKey, nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Short Key not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Short Key not found'. Got '%s'", m["error"])
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
	_, err := App.DB.Exec("INSERT INTO urls(orignal_url, short_key, expire_time) VALUES($1, $2, $3)", originalURL, shortKey, expireTime)
	return shortKey, err
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
