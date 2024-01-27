package app

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"time"

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

func validateURL(originalURL string) bool {
	if _, err := url.ParseRequestURI(originalURL); err != nil {
		return false
	}
	return true
}

func validateShortKey(shortKey string) bool {
	regexPattern := `^[A-Z a-z 0-9]{10,11}$`
	if ok, _ := regexp.MatchString(regexPattern, shortKey); !ok {
		return false
	}
	return true
}

func setRedisKey(r *redis.Client, ctx context.Context, key string, value interface{}, ttl time.Duration) {
	var err error
	var val []byte
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

func deleteRedisKey(r *redis.Client, ctx context.Context, keys ...string) {
	if _, err := r.Del(ctx, keys...).Result(); err != nil {
		log.Println("[Error] Could not delete key in redis ", err)
	}
}
