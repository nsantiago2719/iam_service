package main

import (
	"fmt"
	"os"
	"strconv"
)

// Getenv retrieves the env value else return a default value
func Getenv(key string, defaultValue ...string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultValue[0]
}

// GetenvInt retrieves env or use default value then returns int
func GetenvInt(key string, defaultValue ...string) int {
	env := Getenv(key, defaultValue[0])

	val, err := strconv.Atoi(env)
	if err != nil {
		fmt.Println("Error: %w", err)
	}

	return val
}

// GetenvBool retrieves env or use default value then returns bool
func GetenvBool(key string, defaultValue ...string) bool {
	env := Getenv(key, defaultValue[0])

	val, err := strconv.ParseBool(env)
	if err != nil {
		fmt.Println("Error: %w", err)
	}

	return val
}
