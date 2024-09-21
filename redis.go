package main

import (
	"errors"

	"github.com/redis/go-redis/v9"
)

func RedisClient() (*redis.Client, error) {
	rcon := redis.NewClient(&redis.Options{
		Addr:             Getenv("REDIS_ADDR"),
		Password:         Getenv("REDIS_PASS", ""),
		DB:               GetenvInt("REDIS_DB", "0"),
		Protocol:         GetenvInt("REDIS_PROTOCOL", "3"),
		DisableIndentity: GetenvBool("REDIS_IDENITY", "false"),
	})

	if rcon == nil {
		return nil, errors.New("Error connecting to Redis")
	}
	return rcon, nil
}
