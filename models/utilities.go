package models

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

const base62Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func stringToBase62(input string) string {
	var result string
	var numericValue int
	base := len(base62Alphabet)

	for i := 0; i < len(input); i++ {
		charIndex := strings.Index(base62Alphabet, string(input[i]))
		numericValue = numericValue*base + charIndex
	}

	for numericValue > 0 {
		remainder := numericValue % base
		result = string(base62Alphabet[remainder]) + result
		numericValue /= base
	}

	return result
}

func calculateMD5(input string) string {
	hashInBytes := md5.Sum([]byte(input))
	hashString := hex.EncodeToString(hashInBytes[:])
	return hashString
}
