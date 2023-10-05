package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type config struct {
	host          string
	port          string
	redisAddr     string
	redisUser     string
	redisPassword string
	redisDB       int
	urlLength     int
}

func getConfig() config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	port := os.Getenv("PLAKKEN_PORT")
	redisAddr := os.Getenv("PLAKKEN_REDIS_ADDR")
	db := os.Getenv("PLAKKEN_REDIS_DB")
	uLen := os.Getenv("PLAKKEN_REDIS_URL_LEN")

	if port == "" || redisAddr == "" {
		log.Fatal("Missing or invalid PLAKKEN_PORT or PLAKKEN_REDIS_ADDR")
	}

	redisDB, err := strconv.Atoi(db)
	if err != nil {
		log.Fatal("Invalid PLAKKEN_REDIS_DB")
	}

	urlLen, err := strconv.Atoi(uLen)
	if err != nil {
		log.Fatal("Invalid PLAKKEN_REDIS_URL_LEN")
	}

	conf := config{
		host:          os.Getenv("PLAKKEN_INTERFACE"),
		port:          port,
		redisAddr:     redisAddr,
		redisUser:     os.Getenv("PLAKKEN_REDIS_USER"),
		redisPassword: os.Getenv("PLAKKEN_REDIS_PASSWORD"),
		redisDB:       redisDB,
		urlLength:     urlLen,
	}

	return conf
}
