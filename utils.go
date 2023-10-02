package main

import "math/rand"

func generateUrl() string {
	length := currentConfig.urlLength
	listChars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = listChars[rand.Intn(len(listChars))]
	}

	return string(b)
}
