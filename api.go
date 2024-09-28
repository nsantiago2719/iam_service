package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type API struct {
	listenAddr  string
	database    *PostgresDb
	memoryCache *redis.Client
}

func APIServer(listenAddr, postgresDsn string) *API {
	db, err := PostgresCreate(postgresDsn)
	if err != nil {
		log.Fatal("Error connecting to database", err)
	}

	memDb, err := RedisClient()
	if err != nil {
		log.Fatal("Error connecting to redis database", err)
	}
	return &API{
		listenAddr:  listenAddr,
		database:    db,
		memoryCache: memDb,
	}
}

func (s *API) Create() {
	router := mux.NewRouter()

	s.database.db.MustExec(schema)
	s.IamRoutes(router)

	log.Println("iam service running on port", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *API) IamRoutes(r *mux.Router) {
	userPrefix := r.PathPrefix("/users").Subrouter()
	r.HandleFunc("/auth", s.Auth).Methods(http.MethodPost)
	r.HandleFunc("/logout", s.Logout).Methods(http.MethodDelete)

	userPrefix.HandleFunc("/{id}", UserDetails).Methods(http.MethodGet)
}
