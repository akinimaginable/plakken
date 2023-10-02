package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

var ctx = context.Background()

func connectDB() *redis.Client {
	db := redis.NewClient(&redis.Options{
		Addr:     currentConfig.redisAddr,
		Username: currentConfig.redisUser,
		Password: currentConfig.redisPassword,
		DB:       currentConfig.redisDB,
	})
	return db
}

func insert_paste(key string, content string, secret string, ttl time.Duration) {
	type dbSchema struct {
		content string
		secret  string
	}

	hash := dbSchema{
		content: content,
		secret:  secret,
	}
	connectDB().HSet(ctx, key, "content", hash.content)
	connectDB().HSet(ctx, key, "secret", hash.secret)
	if ttl > -1 {
		connectDB().Do(ctx, key, ttl)
	}
}
