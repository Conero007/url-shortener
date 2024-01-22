package models

import (
	"database/sql"
	"time"
)

type URL struct {
	ShortKey    string    `json:"-"`
	OriginalURL string    `json:"original_url"`
	ShortURL    string    `json:"short_url"`
	ExpireTime  time.Time `json:"expire_time"`
}

func (u *URL) Create(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO urls(original_url, short_key, expire_time) VALUES(?, ?, ?)", u.OriginalURL, u.ShortKey, u.ExpireTime)
	return err
}

func (u *URL) Fetch(db *sql.DB) {
	db.QueryRow("SELECT original_url, short_key, expire_time FROM urls WHERE short_key = ?", u.ShortKey).Scan(&u.OriginalURL, &u.ShortKey, &u.ExpireTime)
}
