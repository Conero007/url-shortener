package main

import (
	"fmt"
	"math/rand"
	"net/http"
)

type URLShortener struct {
	urls map[string]string
}

func MakeURLShortener() *URLShortener {
	return &URLShortener{
		urls: make(map[string]string),
	}
}

func (us URLShortener) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	shortKey := r.URL.Path[len("/short/"):]
	if shortKey == "" {
		http.Error(w, "shortened Key is missing", http.StatusBadRequest)
		return
	}

	originalURL, ok := us.urls[shortKey]
	if !ok {
		http.Error(w, "shortened key not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

func (us *URLShortener) HandleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "url parameter missing", http.StatusBadRequest)
		return
	}

	shortKey := generateShortKey()
	us.urls[shortKey] = originalURL

	shortenedURL := fmt.Sprintf("http://localhost:8080/short/%s", shortKey)

	w.Header().Set("Content-Type", "text/html")

	fmt.Fprint(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>URL Shortener</title>
		</head>
		<body>
			<h2>URL Shortener</h2>
			<p>Original URL: `, originalURL, `</p>
			<p>Shortened URL: <a href="`, shortenedURL, `">`, shortenedURL, `</a></p>
		</body>
		</html>
	`)
}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	shortKey := make([]byte, keyLength)

	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}

	return string(shortKey)
}
