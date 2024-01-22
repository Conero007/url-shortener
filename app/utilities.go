package app

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"
)

const keyLength = 6
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateShortKey() string {
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}

func generateShortURL(shortKey string) string {
	return fmt.Sprintf("http://%s:%s/%s", os.Getenv("APP_URL"), os.Getenv("PORT"), shortKey)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func validateURL(originalURL string) bool {
	if _, err := url.ParseRequestURI(originalURL); err != nil {
		return false
	}
	return true
}

func validateShortKey(shortKey string) bool {
	regexPattern := `^[A-Z a-z 0-9]{6}$`
	if ok, _ := regexp.MatchString(regexPattern, shortKey); !ok {
		return false
	}
	return true
}

func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
