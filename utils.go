package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
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

// RemoveStringDuplicate accepts map of string and returns compact string
func RemoveStringDuplicate(s []string) string {
	uniqueMap := map[string]bool{}
	stringMap := []string{}

	for _, v := range s {
		splits := strings.Split(v, " ")
		for _, s := range splits {
			// check if the string is in the unique map
			// if returns true do not append to stringMap
			if !uniqueMap[s] {
				// set the string as key and value as true
				uniqueMap[s] = true
				// append to stringMap
				stringMap = append(stringMap, s)
			}
		}

	}

	return strings.Join(stringMap, " ")
}

func JSONWriter(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
