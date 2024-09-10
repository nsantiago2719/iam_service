package main

import (
	_ "database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	JwtSecret  = []byte(os.Getenv("JWT_SECRET"))
	pgPassword = os.Getenv("PG_PASSWORD")
	pgHost     = os.Getenv("PG_HOST")
	pgUser     = os.Getenv("PG_USER")
	pgPort     = os.Getenv("PG_PORT")
	dsn        = fmt.Sprintf("user=%s password=%s dbname=iam host=%s port=%s sslmode=disable", pgUser, pgPassword, pgHost, pgPort)
)

// Setup DB
var db, err = sqlx.Connect("postgres", dsn)

// Bcrypt Constants
const (
	MinCost     int = 8
	MaxCost     int = 30
	DefaultCost int = 10
)

func main() {
	r := mux.NewRouter()

	db.MustExec(schema)

	IamRoutes(r)

	fmt.Println("Server running and listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
