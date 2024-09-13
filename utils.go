package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

// ExtractClaim returns the jwt claim and bool for token validitiy
func ExtractClaim(token string) (*Claims, bool) {
	tokenVal, err := jwt.ParseWithClaims(token,
		&Claims{},
		func(_ *jwt.Token) (interface{}, error) {
			return JwtSecret, nil
		})
	if err != nil {
		fmt.Println(err)
		return nil, false
	}

	if claims, ok := tokenVal.Claims.(*Claims); ok && tokenVal.Valid {
		ctx := context.Background()
		rclient := RedisClient()
		val, _ := rclient.Get(ctx, claims.RegisteredClaims.ID).Result()

		// check if there is a value then return invalid
		if len(val) > 0 {
			return nil, false
		}

		return claims, true
	}

	// by default return invalid token
	return nil, false
}

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
