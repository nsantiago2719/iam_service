package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func IamRoutes(r *mux.Router) {
	userPrefix := r.PathPrefix("/users").Subrouter()
	r.HandleFunc("/auth", Auth).Methods(http.MethodPost)
	r.HandleFunc("/logout", Logout).Methods(http.MethodDelete)

	userPrefix.HandleFunc("/{id}", UserDetails).Methods(http.MethodGet)
}
