package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Conero007/url-shortener/constants"
	"github.com/go-sql-driver/mysql"
)

type ShortenURL struct {
	ID          int64     `json:"-"`
	ShortKey    string    `json:"-"`
	OriginalURL string    `json:"original_url"`
	ShortURL    string    `json:"short_url"`
	ExpireTime  time.Time `json:"expire_time"`
}

func GetShortenURL(originalURL string) *ShortenURL {
	return &ShortenURL{
		OriginalURL: originalURL,
	}
}

func (u *ShortenURL) CreateShortURL(db *sql.DB) error {
	customShortKey := true
	if u.ShortKey == "" {
		customShortKey = false
		u.generateShortKey(false)
	}

	if u.ExpireTime.IsZero() {
		u.updateExpireTime()
	}

	var attemptCounter int
	var mysqlErr *mysql.MySQLError
	query := "INSERT INTO urls(original_url, short_key, expire_time) VALUES(?, ?, ?);"

	_, err := db.Exec(query, u.OriginalURL, u.ShortKey, u.ExpireTime)

	for !customShortKey && err != nil && errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 && attemptCounter < constants.GENERATE_SHORT_KEY_MAX_ATTEMPT {
		attemptCounter++
		u.generateShortKey(true)
		_, err = db.Exec(query, u.OriginalURL, u.ShortKey, u.ExpireTime)
	}

	u.generateShortURL()

	return err
}

func (u *ShortenURL) FetchShortURLData(db *sql.DB) {
	query := "SELECT id, original_url, short_key, expire_time FROM urls WHERE short_key = ? LIMIT 1;"
	db.QueryRow(query, u.ShortKey).Scan(&u.ID, &u.OriginalURL, &u.ShortKey, &u.ExpireTime)
}

func (u *ShortenURL) DeleteShortURLData(db *sql.DB, wg *sync.WaitGroup) {
	defer wg.Done()

	query := "DELETE FROM urls WHERE id = ? LIMIT 1;"
	if _, err := db.Exec(query, u.ID); err != nil {
		log.Printf("[Error] Could not Delete row %d. ERROR: %s", u.ID, err.Error())
	}
}

func (u *ShortenURL) generateShortKey(retry bool) {
	var md5Result string
	if retry {
		md5Result = calculateMD5(u.OriginalURL + generateRandomKey())
	} else {
		md5Result = calculateMD5(u.OriginalURL)
	}
	base62Result := stringToBase62(md5Result)
	u.ShortKey = base62Result[:constants.SHORT_KEY_LENGTH]
}

func (u *ShortenURL) generateShortURL() {
	u.ShortURL = fmt.Sprintf("http://%s:%s/%s", os.Getenv("APP_URL"), os.Getenv("PORT"), u.ShortKey)
}

func (u *ShortenURL) updateExpireTime() {
	t := FetchMaxExpireTime()
	u.ExpireTime = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
