package utils

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	mathrand "math/rand"
	"strconv"
	"strings"
)

// GenerateUrl generate random string for plak url
func GenerateUrl(length uint8) string {
	listChars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = listChars[mathrand.Intn(len(listChars))]
	}

	return string(b)
}

// GenerateSecret generate random secret (32 bytes hexadecimal)
func GenerateSecret() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		log.Printf("Failed to generate secret")
	}

	return hex.EncodeToString(key)
}

// CheckCharRedundant  verify is a character is redundant in a string
func CheckCharRedundant(source string, char string) bool { // Verify if a char is redundant
	return strings.Count(source, char) > 1
}

func parseIntBeforeSeparator(source *string, sep string) (int, error) { // return 0 & error if error, only accept positive number
	if CheckCharRedundant(*source, sep) {
		return 0, &ParseIntBeforeSeparatorError{Message: *source + ": cannot parse value as int"}
	}
	var value int
	var err error
	if strings.Contains(*source, sep) {
		value, err = strconv.Atoi(strings.Split(*source, sep)[0])
		if err != nil {
			log.Println(err)
			return 0, &ParseIntBeforeSeparatorError{Message: *source + ": cannot parse value as int"}
		}
		if value < 0 { // Only positive value is correct
			return 0, &ParseIntBeforeSeparatorError{Message: *source + ": format only take positive value"}
		}

		if value > 99 {
			return 0, &ParseIntBeforeSeparatorError{Message: *source + ": Format only take two number"}
		}

		*source = strings.Join(strings.Split(*source, sep)[1:], "")
	}
	return value, nil
}

// ParseExpiration Parse "1d1h1m1s" duration format. Return 0 & error if error
func ParseExpiration(source string) (int, error) {
	var expiration int
	var tempOutput int
	var err error
	if source == "0" {
		return 0, nil
	}

	source = strings.ToLower(source)

	tempOutput, err = parseIntBeforeSeparator(&source, "d")
	expiration = tempOutput * 86400
	if err != nil {
		log.Println(err)
		return 0, &ParseExpirationError{Message: "Invalid syntax"}
	}
	tempOutput, err = parseIntBeforeSeparator(&source, "h")
	expiration += tempOutput * 3600
	if err != nil {
		log.Println(err)
		return 0, &ParseExpirationError{Message: "Invalid syntax"}
	}
	tempOutput, err = parseIntBeforeSeparator(&source, "m")
	expiration += tempOutput * 60
	if err != nil {
		log.Println(err)
		return 0, &ParseExpirationError{Message: "Invalid syntax"}
	}
	tempOutput, err = parseIntBeforeSeparator(&source, "s")
	expiration += tempOutput * 1
	if err != nil {
		log.Println(err)
		return 0, &ParseExpirationError{Message: "Invalid syntax"}
	}

	return expiration, nil
}
