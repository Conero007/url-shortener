package app

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Conero007/url-shortener-golang/models"
	"github.com/gorilla/mux"
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

	u := models.GetURL(originalURL)
	if err := u.CreateShortURL(App.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong. Please try again.")
		return
	}

	respondWithJSON(w, http.StatusCreated, &u)
}

func HandleRedirectToOriginalURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if !validateShortKey(vars["key"]) {
		respondWithError(w, http.StatusBadRequest, "Invalid short key")
		return
	}

	u := models.URL{ShortKey: vars["key"]}
	u.Fetch(App.DB)

	if u.OriginalURL == "" || u.ExpireTime.Before(time.Now()) {
		App.wg.Add(1)
		go u.Delete(App.DB, App.wg)
		respondWithError(w, http.StatusNotFound, "Short Key not found")
		return
	}

	http.Redirect(w, r, u.OriginalURL, http.StatusMovedPermanently)
}
