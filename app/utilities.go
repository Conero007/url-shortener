package app

import (
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)

	if App.debug {
		App.wg.Wait()
	}
}

func validateURL(originalURL string) bool {
	if _, err := url.ParseRequestURI(originalURL); err != nil {
		return false
	}
	return true
}

func validateShortKey(shortKey string) bool {
	regexPattern := `^[A-Z a-z 0-9]{10,11}$`
	if ok, _ := regexp.MatchString(regexPattern, shortKey); !ok {
		return false
	}
	return true
}
