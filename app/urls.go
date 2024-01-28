package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Conero007/url-shortener-golang/models"
	"github.com/gorilla/mux"
)

type ShortenURLRequest struct {
	URL            string `json:"url"`
	CustomShortKey string `json:"custom_short_key"`
}

func HandleURLShortening(w http.ResponseWriter, r *http.Request) {
	var requestBody ShortenURLRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if requestBody.URL == "" {
		respondWithError(w, http.StatusBadRequest, "URL not given")
		return
	}

	if !validateURL(requestBody.URL) {
		respondWithError(w, http.StatusBadRequest, "Invalid URL given")
		return
	}

	u := models.GetShortenURL(requestBody.URL)

	if requestBody.CustomShortKey != "" && !validateShortKey(requestBody.CustomShortKey) {
		respondWithError(w, http.StatusBadRequest, "Invalid custom short key")
		return
	} else if requestBody.CustomShortKey != "" && !models.CheckShortKeyAvailability(App.DB, requestBody.CustomShortKey) {
		respondWithError(w, http.StatusNotAcceptable, "Short key not available to use")
		return
	}

	u.ShortKey = requestBody.CustomShortKey

	if err := u.CreateShortURL(App.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong. Please try again.")
		return
	}

	App.wg.Add(1)
	go setRedisKey(App.Redis, context.Background(), App.wg, u.ShortKey, u, 24*time.Hour)

	respondWithJSON(w, http.StatusCreated, &u)
}

func HandleRedirectToOriginalURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if !validateShortKey(vars["key"]) {
		respondWithError(w, http.StatusBadRequest, "Invalid short key")
		return
	}

	var u models.ShortenURL

	if err := getRedisKey(App.Redis, context.Background(), vars["key"], &u); err != nil {
		u.ShortKey = vars["key"]
		u.FetchShortURLData(App.DB)
	}

	if u.OriginalURL == "" || u.ExpireTime.Before(time.Now()) {
		App.wg.Add(2)
		go u.DeleteShortURLData(App.DB, App.wg)
		go deleteRedisKey(App.Redis, context.Background(), App.wg, vars["key"])
		respondWithError(w, http.StatusNotFound, "Short Key not found")
		return
	}

	http.Redirect(w, r, u.OriginalURL, http.StatusMovedPermanently)
}
