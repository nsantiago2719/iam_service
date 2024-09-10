package main

import (
	"fmt"
	"os"
	"strconv"
)

func Getenv(key string, defaultValue ...string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultValue[0]
}

func GetenvInt(key string, defaultValue ...string) int {
	env := Getenv(key, defaultValue[0])

	val, err := strconv.Atoi(env)
	if err != nil {
		fmt.Println("Error: %w", err)
	}

	return val
}

func GetenvBool(key string, defaultValue ...string) bool {
	env := Getenv(key, defaultValue[0])

	val, err := strconv.ParseBool(env)
	if err != nil {
		fmt.Println("Error: %w", err)
	}

	return val
}
