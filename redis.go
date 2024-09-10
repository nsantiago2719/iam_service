package main

import (
	"github.com/redis/go-redis/v9"
)

func RedisClient() *redis.Client {
	rcon := redis.NewClient(&redis.Options{
		Addr:     Getenv("REDIS_ADDR"),
		Password: Getenv("REDIS_PASS"),
		DB:       GetenvInt("REDIS_DB"),
		Protocol: GetenvInt("REDIS_PROTOCOL"),
	})
}
