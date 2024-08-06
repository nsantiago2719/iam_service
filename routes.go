package main

import (
	"github.com/gorilla/mux"
)

func IamRoutes(r *mux.Router) {
	r.HandleFunc("/auth", Auth).Methods("POST")
	r.HandleFunc("/user/create", CreateUser).Methods("POST")
}
