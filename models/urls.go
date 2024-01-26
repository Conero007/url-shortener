package models

import (
	"database/sql"
	"log"
	"sync"
	"time"
)

type URL struct {
	ID          int64     `json:"-"`
	ShortKey    string    `json:"-"`
	OriginalURL string    `json:"original_url"`
	ShortURL    string    `json:"short_url"`
	ExpireTime  time.Time `json:"expire_time"`
}

func (u *URL) Create(db *sql.DB, wg *sync.WaitGroup) {
	defer wg.Done()

	if row, err := db.Exec("INSERT INTO urls(original_url, short_key, expire_time) VALUES(?, ?, ?)", u.OriginalURL, u.ShortKey, u.ExpireTime); err != nil {
		log.Printf("[Error] Could insert into DB. ERROR: %s", err.Error())
	} else if u.ID, err = row.LastInsertId(); err != nil {
		log.Printf("[Error] Could fetch ID of the inserted row. ERROR: %s", err.Error())
	}
}

func (u *URL) Fetch(db *sql.DB) {
	db.QueryRow("SELECT id, original_url, short_key, expire_time FROM urls WHERE short_key = ?", u.ShortKey).Scan(&u.ID, &u.OriginalURL, &u.ShortKey, &u.ExpireTime)
}

func (u *URL) Delete(db *sql.DB, wg *sync.WaitGroup) {
	defer wg.Done()

	if _, err := db.Exec("DELETE FROM urls WHERE id = ? LIMIT 1", u.ID); err != nil {
		log.Printf("[Error] Could not Delete row %d. ERROR: %s", u.ID, err.Error())
	}
}
