package app

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Conero007/url-shortener-golang/models"
)

func HandleURLShortening(w http.ResponseWriter, r *http.Request) {
	var requestBody map[string]interface{}
	json.NewDecoder(r.Body).Decode(&requestBody)

	_, ok := requestBody["url"]
	if !ok {
		respondWithError(w, http.StatusBadRequest, "URL not given")
		return
	}

	originalURL := requestBody["url"].(string)

	if !validateURL(originalURL) {
		respondWithError(w, http.StatusBadRequest, "Invalid URL given")
		return
	}

	shortKey := generateShortKey()
	shortURL := generateShortURL(shortKey)
	expireTime := truncateToDay(time.Now().AddDate(0, 0, 8))

	u := models.URL{
		OriginalURL: originalURL,
		ShortKey:    shortKey,
		ShortURL:    shortURL,
		ExpireTime:  expireTime,
	}

	if err := u.Create(App.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		log.Print(err)
		return
	}

	respondWithJSON(w, http.StatusCreated, u)
}

func HandleRedirectToOriginalURL(w http.ResponseWriter, r *http.Request) {
}
