package main

import (
	"github.com/redis/go-redis/v9"
)

func RedisClient() *redis.Client {
	rcon := redis.NewClient(&redis.Options{
		Addr:             Getenv("REDIS_ADDR"),
		Password:         Getenv("REDIS_PASS", ""),
		DB:               GetenvInt("REDIS_DB", "0"),
		Protocol:         GetenvInt("REDIS_PROTOCOL", "3"),
		DisableIndentity: GetenvBool("REDIS_IDENITY", "false"),
	})

	return rcon
}
