package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func IamRoutes(r *mux.Router) {
	r.HandleFunc("/auth", auth).Methods("GET")
}

func auth(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("")
}
