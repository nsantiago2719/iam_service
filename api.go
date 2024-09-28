package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

func (e APIError) Error() string {
	return e.Msg
}

func JSONWriter(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeHandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			if e, ok := err.(APIError); ok {
				slog.Error(e.Path, "msg", e.Msg, "status", e.Status)
				response := GenericResponse{
					Message: err.Error(),
				}
				JSONWriter(w, http.StatusUnauthorized, response)
			}
		}
	}
}

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
	r.HandleFunc("/auth", makeHandler(s.handleAuth)).Methods(http.MethodPost)
	r.HandleFunc("/logout", makeHandler(s.handleLogout)).Methods(http.MethodDelete)

	userPrefix.HandleFunc("/{id}", UserDetails).Methods(http.MethodGet)
}
