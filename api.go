package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type API struct {
	listenAddr string
	database   *sqlx.DB
}

func APIServer(listenAddr, postgresDsn string) *API {
	db, err := PostgresCreate(postgresDsn)
	if err != nil {
		log.Fatal("Error connecting to database", err)
	}
	return &API{
		listenAddr: listenAddr,
		database:   db,
	}
}

func (s *API) Create() {
	router := mux.NewRouter()

	s.database.MustExec(schema)
	s.IamRoutes(router)

	log.Println("iam service running on port", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func JSONWriter(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
