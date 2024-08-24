package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func UserDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	parameters := mux.Vars(r)

	response := map[string]string{
		"id": parameters["id"],
	}

	json.NewEncoder(w).Encode(response)
}
