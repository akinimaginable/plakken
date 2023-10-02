package main

import (
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
}

func setConfig() config {
	host := os.Getenv("PLAKKEN_HOST")

	port := os.Getenv("PLAKKEN_PORT")
	if port == "" {
		port = "3000"
	}
	redisAddr := os.Getenv("PLAKKEN_REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisUser := os.Getenv("PLAKKEN_REDIS_USER")
	redisPassword := os.Getenv("PLAKKEN_REDIS_PASSWORD")
	redisDB, err := strconv.Atoi(os.Getenv("PLAKKEN_REDIS_DB"))
	if err != nil {
		redisDB = 0
	}

	s := config{
		host:          host,
		port:          port,
		redisAddr:     redisAddr,
		redisUser:     redisUser,
		redisPassword: redisPassword,
		redisDB:       redisDB,
	}

	return s
}
