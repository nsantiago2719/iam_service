package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *API) IamRoutes(r *mux.Router) {
	userPrefix := r.PathPrefix("/users").Subrouter()
	r.HandleFunc("/auth", s.Auth).Methods(http.MethodPost)
	r.HandleFunc("/logout", s.Logout).Methods(http.MethodDelete)

	userPrefix.HandleFunc("/{id}", UserDetails).Methods(http.MethodGet)
}
