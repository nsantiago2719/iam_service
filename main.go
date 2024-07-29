package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nsantiago2719/i_a_m/config"
)

func main() {
	r := mux.NewRouter()

	config.IamRoutes(r)

	fmt.Println("Server running and listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
