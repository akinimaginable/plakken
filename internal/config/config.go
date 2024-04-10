package config

import (
	"log"
	"os"
	"strconv"

	"git.gnous.eu/gnouseu/plakken/internal/constant"
)

// InitConfig Structure for program initialisation.
type InitConfig struct {
	ListenAddress string
	RedisAddress  string
	RedisUser     string
	RedisPassword string
	RedisDB       int
	URLLength     uint8
}

// GetConfig Initialise configuration form .env.
func GetConfig() InitConfig {
	listenAddress := os.Getenv("PLAKKEN_LISTEN")
	redisAddress := os.Getenv("PLAKKEN_REDIS_ADDRESS")
	db := os.Getenv("PLAKKEN_REDIS_DB")
	uLen := os.Getenv("PLAKKEN_URL_LENGTH")

	if listenAddress == "" || redisAddress == "" {
		log.Fatal("Missing or invalid listenAddress or PLAKKEN_REDIS_ADDRESS")
	}

	redisDB, err := strconv.Atoi(db)
	if err != nil {
		log.Fatal("Invalid PLAKKEN_REDIS_DB")
	}

	urlLength, err := strconv.Atoi(uLen)
	if err != nil {
		log.Fatal("Invalid PLAKKEN_URL_LENGTH")
	}

	if urlLength > constant.MaxURLLength {
		log.Fatal("PLAKKEN_URL_LENGTH cannot be greater than 255")
	}

	return InitConfig{
		ListenAddress: listenAddress,
		RedisAddress:  redisAddress,
		RedisUser:     os.Getenv("PLAKKEN_REDIS_USER"),
		RedisPassword: os.Getenv("PLAKKEN_REDIS_PASSWORD"),
		RedisDB:       redisDB,
		URLLength:     uint8(urlLength),
	}
}
