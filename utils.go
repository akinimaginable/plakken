package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	mathrand "math/rand"
	"strconv"
	"strings"
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

func parseIntBeforeSeparator(source *string, sep string) (int, error) { // return -1 if error, only accept positive number
	var value int
	var err error
	if strings.Contains(*source, sep) {
		value, err = strconv.Atoi(strings.Split(*source, sep)[0])
		if err != nil {
			log.Println(err)
			return -1, fmt.Errorf("parseIntBeforeSeparator : \"%s\" : cannot parse value as int", *source)
		}
		if value < 0 { // Only positive value is correct
			return -1, fmt.Errorf("parseIntBeforeSeparator : \"%s\" : format only take positive value", *source)
		}
		*source = strings.Join(strings.Split(*source, sep)[1:], "")
	}
	return value, nil
}

func ParseExpiration(source string) (int, error) { // return -1 if error
	var expiration int
	var tempOutput int
	var err error
	errMessage := "ParseExpiration : \"%s\" : invalid syntax"
	if source == "0" {
		return 0, nil
	}
	tempOutput, err = parseIntBeforeSeparator(&source, "d")
	expiration = tempOutput * 86400
	if err != nil {
		log.Println(err)
		return -1, fmt.Errorf(errMessage, source)
	}
	tempOutput, err = parseIntBeforeSeparator(&source, "h")
	expiration += tempOutput * 3600
	if err != nil {
		log.Println(err)
		return -1, fmt.Errorf(errMessage, source)
	}
	tempOutput, err = parseIntBeforeSeparator(&source, "m")
	expiration += tempOutput * 60
	if err != nil {
		log.Println(err)
		return -1, fmt.Errorf(errMessage, source)
	}
	tempOutput, err = parseIntBeforeSeparator(&source, "s")
	expiration += tempOutput * 1
	if err != nil {
		log.Println(err)
		return -1, fmt.Errorf(errMessage, source)
	}

	return expiration, nil
}
