package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

type URL struct {
	ID          int64     `json:"-"`
	ShortKey    string    `json:"-"`
	OriginalURL string    `json:"original_url"`
	ShortURL    string    `json:"short_url"`
	ExpireTime  time.Time `json:"expire_time"`
}

func GetURL(originalURL string) *URL {
	return &URL{
		OriginalURL: originalURL,
	}
}

func (u *URL) CreateShortURL(db *sql.DB) error {
	u.generateShortKey()
	u.generateShortURL()
	u.updateExpireTime()

	var e *pgconn.PgError
	_, err := db.Exec("INSERT INTO urls(original_url, short_key, expire_time) VALUES(?, ?, ?)", u.OriginalURL, u.ShortKey, u.ExpireTime)
	if err != nil && errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
		log.Printf("[Error] Could insert into DB. ERROR: %s", err.Error())
		return err
	}

	return nil
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

func (u *URL) generateShortKey() {
	md5Result := calculateMD5(u.OriginalURL)
	base62Result := stringToBase62(md5Result)
	u.ShortKey = base62Result
}

func (u *URL) generateShortURL() {
	u.ShortURL = fmt.Sprintf("http://%s:%s/%s", os.Getenv("APP_URL"), os.Getenv("PORT"), u.ShortKey)
}

func (u *URL) updateExpireTime() {
	t := time.Now().AddDate(0, 0, 8)
	u.ExpireTime = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
