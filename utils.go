package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	mathrand "math/rand"
)

func generateUrl() string {
	listChars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, currentConfig.urlLength)
	for i := range b {
		b[i] = listChars[mathrand.Intn(len(listChars))]
	}

	return string(b)
}

func generateSecret() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		log.Printf("Failed to generate secret")
	}

	return hex.EncodeToString(key)
}

func urlExist(url string) bool {
	exist := db.Exists(ctx, url).Val()
	return exist == 1
}

func verifySecret(url string, secret string) bool {
	if secret == db.HGet(ctx, url, "secret").Val() {
		return true
	}
	return false
}
