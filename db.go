package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresDb struct {
	db *sqlx.DB
}

func PostgresCreate(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("Error connection to the database:", err)
	}
	return db, nil
}
