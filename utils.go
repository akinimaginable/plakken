package main

import (
	"crypto/rand"
	"encoding/hex"
	mathrand "math/rand"
)

func generateUrl() string {
	length := currentConfig.urlLength
	listChars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = listChars[mathrand.Intn(len(listChars))]
	}

	return string(b)
}

func generateSecret() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		// handle error here
	}

	secret := hex.EncodeToString(key)
	return secret
}

func urlExist(url string) bool {
	exist := connectDB().Exists(ctx, url).Val()
	if exist == 1 {
		return true
	}
	return false
}
