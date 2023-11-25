package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	mathrand "math/rand"
)

func GenerateUrl() string {
	listChars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, currentConfig.urlLength)
	for i := range b {
		b[i] = listChars[mathrand.Intn(len(listChars))]
	}

	return string(b)
}

func GenerateSecret() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		log.Printf("Failed to generate secret")
	}

	return hex.EncodeToString(key)
}

func UrlExist(url string) bool {
	return db.Exists(ctx, url).Val() == 1
}

func VerifySecret(url string, secret string) bool {
	return secret == db.HGet(ctx, url, "secret").Val()
}
