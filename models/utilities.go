package models

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"math/rand"
	"strings"
	"time"

	"github.com/Conero007/url-shortener-golang/constants"
)

func CheckShortKeyAvailability(db *sql.DB, customShortKey string) bool {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM urls WHERE short_key = ? LIMIT 1)"
	db.QueryRow(query, customShortKey).Scan(&exists)
	return exists
}

func stringToBase62(input string) string {
	var result string
	var numericValue uint64
	base := len(constants.BASE_62_CHARACTERS)

	for i := 0; i < len(input); i++ {
		charIndex := strings.Index(constants.BASE_62_CHARACTERS, string(input[i]))
		numericValue = numericValue*uint64(base) + uint64(charIndex)
	}

	for numericValue > 0 {
		remainder := numericValue % uint64(base)
		result = string(constants.BASE_62_CHARACTERS[remainder]) + result
		numericValue /= uint64(base)
	}

	return result
}

func calculateMD5(input string) string {
	hashInBytes := md5.Sum([]byte(input))
	hashString := hex.EncodeToString(hashInBytes[:])
	return hashString
}

func generateRandomKey() string {
	var randomString []byte
	hash := calculateMD5(time.Now().String())
	for len(randomString) < constants.RANDOM_KEY_LENGTH {
		randomString = append(randomString, hash[rand.Intn(len(hash))])
	}
	return string(randomString)
}

func FetchMaxExpireTime() time.Time {
	return time.Now().AddDate(0, 0, 8)
}
