package database

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type DBConfig struct {
	DB *redis.Client
}

var ctx = context.Background()

// InitDB initialise redis connection settings
func InitDB(addr string, user string, password string, db int) *redis.Options {
	DBConfig := &redis.Options{
		Addr:     addr,
		Username: user,
		Password: password,
		DB:       db,
	}

	return DBConfig
}

// ConnectDB make new database connection
func ConnectDB(dbConfig *redis.Options) *redis.Client {
	localDb := redis.NewClient(dbConfig)
	return localDb
}

// Ping test connection to Redis database
func Ping(db *redis.Client) error {
	status := db.Ping(ctx)
	if status.String() != "ping: PONG" {
		return &PingError{}
	}
	return nil
}

func (config DBConfig) InsertPaste(key string, content string, secret string, ttl time.Duration) {
	type dbSchema struct {
		content string
		secret  string
	}

	hash := dbSchema{
		content: content,
		secret:  secret,
	}
	err := config.DB.HSet(ctx, key, "content", hash.content).Err()
	if err != nil {
		log.Println(err)
	}
	err = config.DB.HSet(ctx, key, "secret", hash.secret).Err()
	if err != nil {
		log.Println(err)
	}
	if ttl > -1 {
		config.DB.Expire(ctx, key, ttl)
	}
}

func (config DBConfig) UrlExist(url string) bool {
	return config.DB.Exists(ctx, url).Val() == 1
}

func (config DBConfig) VerifySecret(url string, secret string) bool {
	return secret == config.DB.HGet(ctx, url, "secret").Val()
}
