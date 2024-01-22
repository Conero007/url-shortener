package main

import (
	"database/sql"
	"errors"
)

type URL struct {
	ID          int    `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortKey    string `json:"short_key"`
}

func (u *URL) getOriginalURL(db *sql.DB) (string, error) {
	return "", errors.New("not implemented")
}

func (u *URL) createShortURL(db *sql.DB) error {
	return errors.New("not implemented")
}
