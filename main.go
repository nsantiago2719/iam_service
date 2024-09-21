package main

import (
	_ "database/sql"
	"fmt"
	"os"
)

var (
	jwtSecret  = []byte(os.Getenv("JWT_SECRET"))
	pgPassword = Getenv("PG_PASSWORD")
	pgHost     = Getenv("PG_HOST")
	pgUser     = Getenv("PG_USER")
	pgPort     = Getenv("PG_PORT")
	dsn        = fmt.Sprintf("user=%s password=%s dbname=iam host=%s port=%s sslmode=disable", pgUser, pgPassword, pgHost, pgPort)
)

// Bcrypt Constants
const (
	MinCost     int = 8
	MaxCost     int = 30
	DefaultCost int = 10
)

func main() {
	fmt.Println(dsn)
	server := APIServer(":8000", dsn)
	server.Create()
}
