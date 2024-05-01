package utils

import (
	"errors"
	"log"
	mathrand "math/rand/v2"
	"os"
	"regexp"
	"strconv"
	"strings"

	"git.gnous.eu/gnouseu/plakken/internal/constant"
)

// GenerateURL generate random string for plak url.
func GenerateURL(length uint8) string {
	listChars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = listChars[mathrand.IntN(len(listChars))]
	}

	return string(b)
}

// CheckCharRedundant  verify is a character is redundant in a string.
func CheckCharRedundant(source string, char string) bool { // Verify if a char is redundant
	return strings.Count(source, char) > 1
}

func parseIntBeforeSeparator(source *string, sep string) (int, error) { // return 0 & error if error, only accept positive number
	if CheckCharRedundant(*source, sep) {
		return 0, &parseIntBeforeSeparatorError{message: *source + ": cannot parse value as int"}
	}
	var value int
	var err error
	if strings.Contains(*source, sep) {
		value, err = strconv.Atoi(strings.Split(*source, sep)[0])
		if err != nil {
			log.Println(err)

			return 0, &parseIntBeforeSeparatorError{message: *source + ": cannot parse value as int"}
		}
		if value < 0 { // Only positive value is correct
			return 0, &parseIntBeforeSeparatorError{message: *source + ": format only take positive value"}
		}

		if value > 99 { //nolint:gomnd
			return 0, &parseIntBeforeSeparatorError{message: *source + ": Format only take two number"}
		}

		*source = strings.Join(strings.Split(*source, sep)[1:], "")
	}

	return value, nil
}

// ParseExpiration Parse "1d1h1m1s" duration format. Return 0 & error if error.
func ParseExpiration(source string) (int, error) {
	var expiration int
	var tempOutput int
	var err error
	if source == "0" {
		return 0, nil
	}

	source = strings.ToLower(source)

	tempOutput, err = parseIntBeforeSeparator(&source, "d")
	expiration = tempOutput * constant.SecondsInDay
	if err != nil {
		log.Println(err)

		return 0, &ParseExpirationError{message: "Invalid syntax"}
	}
	tempOutput, err = parseIntBeforeSeparator(&source, "h")
	expiration += tempOutput * constant.SecondsInHour
	if err != nil {
		log.Println(err)

		return 0, &ParseExpirationError{message: "Invalid syntax"}
	}
	tempOutput, err = parseIntBeforeSeparator(&source, "m")
	expiration += tempOutput * constant.SecondsInMinute
	if err != nil {
		log.Println(err)

		return 0, &ParseExpirationError{message: "Invalid syntax"}
	}
	tempOutput, err = parseIntBeforeSeparator(&source, "s")
	expiration += tempOutput * 1
	if err != nil {
		log.Println(err)

		return 0, &ParseExpirationError{message: "Invalid syntax"}
	}

	return expiration, nil
}

// ValidKey Verify if a key is valid (only letter, number, - and _).
func ValidKey(key string) bool {
	result, err := regexp.MatchString("^[a-zA-Z0-9_.-]*$", key)
	if err != nil {
		return false
	}
	log.Println(key, result)

	return result
}

// FileExist verify if a file exist.
func FileExist(path string) bool {
	_, err := os.Stat(path)

	return !errors.Is(err, os.ErrNotExist)
}
