package handlers

import (
	"encoding/json"
	"github.com/nsantiago2719/i_a_m/models"
	"net/http"
)

func Auth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := ""
	json.NewEncoder(w).Encode(user)
}
