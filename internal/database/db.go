package database

import (
	"context"
	"log"
	"time"

	"git.gnous.eu/gnouseu/plakken/internal/secret"
	"github.com/redis/go-redis/v9"
)

type DBConfig struct {
	DB *redis.Client
}

var ctx = context.Background() //nolint:gochecknoglobals

// InitDB initialise redis connection settings.
func InitDB(addr string, user string, password string, db int) *redis.Options {
	DBConfig := &redis.Options{
		Addr:     addr,
		Username: user,
		Password: password,
		DB:       db,
	}

	return DBConfig
}

// ConnectDB make new database connection.
func ConnectDB(dbConfig *redis.Options) *redis.Client {
	localDB := redis.NewClient(dbConfig)

	return localDB
}

// Ping test connection to Redis database.
func Ping(db *redis.Client) error {
	status := db.Ping(ctx)
	if status.String() != "ping: PONG" {
		return &pingError{}
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

func (config DBConfig) URLExist(url string) bool {
	return config.DB.Exists(ctx, url).Val() == 1
}

func (config DBConfig) VerifySecret(url string, token string) (bool, error) {
	storedHash := config.DB.HGet(ctx, url, "secret").Val()

	result, err := secret.VerifyPassword(token, storedHash)
	if err != nil {
		return false, err
	}

	return result, nil
}
