package config

import (
	"github.com/gorilla/mux"
	"github.com/nsantiago2719/i_a_m/handlers"
)

func IamRoutes(r *mux.Router) {
	r.HandleFunc("/auth", handlers.Auth).Methods("GET")
}
