package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func IamRoutes(r *mux.Router) {
	userPrefix := r.PathPrefix("/users").Subrouter()
	r.HandleFunc("/auth", Auth).Methods(http.MethodPost)

	userPrefix.HandleFunc("/{id}", UserDetails).Methods(http.MethodGet)
}
