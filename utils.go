package main

import (
	"fmt"
	"os"
	"strconv"
)

func Getenv(key string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return ""
}

func GetenvInt(key string) int {
	env := Getenv(key)

	val, err := strconv.Atoi(env)
	if err != nil {
		fmt.Println("Error: %w", err)
	}

	return val
}

func GetenvBool(key string) bool {
	env := Getenv(key)

	val, err := strconv.ParseBool(env)
	if err != nil {
		fmt.Println("Error: %w", err)
	}

	return val
}
