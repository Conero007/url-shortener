package app

import (
	"net/http"
)

func InitializeRoutes() {
	App.Router.HandleFunc("/shorten", HandleURLShortening).Methods(http.MethodPost)
	App.Router.HandleFunc("/{key}", HandleRedirectToOriginalURL).Methods(http.MethodGet)
}
