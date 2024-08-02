package handlers

import (
	"encoding/json"
	"github.com/nsantiago2719/i_a_m/models"
	"net/http"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	permission := models.Permission{
		Action:   "GET",
		Resource: "*",
		Access:   "allow",
	}

	role := models.Role{
		Name: "admin",
		Permissions: []models.Permission{
			permission,
		},
	}

	user := models.User{
		FirstName: "John",
		LastName:  "Doe",
		BirthDate: "11/11/1991",
		Password:  "123123123",
		Username:  "JohnDoe01",
		Email:     "john@doe.com",
		UserRole:  role,
	}

	json.NewEncoder(w).Encode(user)
}
