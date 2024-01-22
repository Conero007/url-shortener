package models

import (
	"database/sql"
	"fmt"
	"time"
)

type URL struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
	ExpireDate  string `json:"expire_date"`

	ShortKey   string    `json:"-"`
	ExpireTime time.Time `json:"-"`
}

func (u *URL) Create(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO urls(original_url, short_key, expire_time) VALUES(?, ?, ?)", u.OriginalURL, u.ShortKey, u.ExpireTime)
	return err
}

func (u *URL) Fetch(db *sql.DB) {
	db.QueryRow("SELECT * FROM urls WHERE short_key = $1", u.ShortKey).Scan(&u)
	fmt.Printf("\n\n%+v\n\n", u)
}
