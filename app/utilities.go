package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"sync"
	"time"

	"github.com/Conero007/url-shortener/constants"
	"github.com/redis/go-redis/v9"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)

	if App.debug {
		App.wg.Wait()
	}
}

func setRedisKey(r *redis.Client, ctx context.Context, wg *sync.WaitGroup, key string, value interface{}, ttl time.Duration) {
	var err error
	var val []byte
	defer wg.Done()

	if val, err = json.Marshal(value); err == nil {
		_, err = r.Set(ctx, key, val, ttl).Result()
	}

	if err != nil {
		log.Println("[Error] Could not set key in redis ", err)
	}
}

func getRedisKey(r *redis.Client, ctx context.Context, key string, dest interface{}) error {
	var err error
	var val string

	if val, err = r.Get(ctx, key).Result(); err == nil {
		return json.Unmarshal([]byte(val), dest)
	}

	if err != redis.Nil {
		log.Println("[Error] Could not get key in redis ", err)
	}

	return err
}

func deleteRedisKey(r *redis.Client, ctx context.Context, wg *sync.WaitGroup, keys ...string) {
	defer wg.Done()

	if _, err := r.Del(ctx, keys...).Result(); err != nil {
		log.Println("[Error] Could not delete key in redis ", err)
	}
}

func validateURL(originalURL string) bool {
	if _, err := url.ParseRequestURI(originalURL); err != nil {
		return false
	}
	return true
}

func validateShortKey(shortKey string) bool {
	regexPattern := fmt.Sprintf(`^[A-Z a-z 0-9]{%d}$`, constants.SHORT_KEY_LENGTH)
	if ok, _ := regexp.MatchString(regexPattern, shortKey); !ok {
		return false
	}
	return true
}
