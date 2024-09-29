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

// JSONWriter is helper to return a json response based on status and v
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

// API details
type API struct {
	listenAddr  string
	database    *PostgresDb
	memoryCache *redis.Client
}

// APIServer initializes the APIServer including
// its db and memDb
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

// Create runs the api
func (s *API) Create() {
	router := mux.NewRouter()

	s.database.db.MustExec(schema)
	s.createRoutes(router)

	log.Println("iam service running on port", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *API) createRoutes(r *mux.Router) *mux.Router {
	userPrefix := r.PathPrefix("/users").Subrouter()
	r.HandleFunc("/auth", makeHandler(s.handleAuth)).Methods(http.MethodPost)
	r.HandleFunc("/logout", makeHandler(s.handleLogout)).Methods(http.MethodDelete)

	userPrefix.HandleFunc("/{id}", UserDetails).Methods(http.MethodGet)
	return r
}
