package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

var ctx = context.Background()

func ConnectDB() *redis.Client {
	localDb := redis.NewClient(&redis.Options{
		Addr:     currentConfig.redisAddr,
		Username: currentConfig.redisUser,
		Password: currentConfig.redisPassword,
		DB:       currentConfig.redisDB,
	})
	return localDb
}

func insertPaste(key string, content string, secret string, ttl time.Duration) {
	type dbSchema struct {
		content string
		secret  string
	}

	hash := dbSchema{
		content: content,
		secret:  secret,
	}
	err := db.HSet(ctx, key, "content", hash.content)
	if err != nil {
		log.Println(err)
	}
	err = db.HSet(ctx, key, "secret", hash.secret)
	if ttl > -1 {
		db.Expire(ctx, key, ttl)
	}
}

func getContent(key string) string {
	return db.HGet(ctx, key, "content").Val()
}

func deleteContent(key string) {
	err := db.Del(ctx, key)
	if err != nil {
		log.Println(err)
	}
}
