package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"regexp"
	"testing"
)

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

	shortKey, err := addShortKey("https://www.google.com/")
	if err != nil {
		t.Errorf("Failed to add short and orign to DB. ERROR: %s", err.Error())
		return
	}

	req, _ := http.NewRequest("GET", "/"+shortKey, nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusMovedPermanently, response.Code)

	if response.Result().Header.Get("Location") != "https://www.google.com/" {
		t.Error("Excepted redirect url https://www.google.com/, found ", response.Result().Header.Get("Location"))
	}
}
