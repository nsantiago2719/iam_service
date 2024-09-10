package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Auth function is used for authentication of the user,
// returns a jwt containing the roles
// and authorized actions
func Auth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	login := LoginDetails{}
	err = json.Unmarshal(body, &login)
	if err != nil {
		fmt.Println("Failed unmarshal of login details: %w", err)
	}

	userRoles := []UserRole{}
	usersMap := make(map[string]*User)
	user := User{}
	query := `
  SELECT users.id AS "users.id",
       users.first_name AS "users.first_name",
       users.last_name AS "users.last_name",
       users.email AS "users.email",
       users.username AS "users.username",
       users.password AS "users.password",
       users.birthdate AS "users.birthdate",
       roles.name AS "roles.name",
       roles.permissions AS "roles.permissions",
       ur.user_id AS "userID"
  FROM users AS users
  LEFT JOIN users_roles AS ur ON users.id = ur.user_id
  LEFT JOIN roles AS roles ON roles.id = ur.role_id
  WHERE users.username=$1
  `

	err = db.Select(&userRoles, query, login.Username)
	if err != nil {
		fmt.Println("User select query failed: %w", err)
	}

	for _, userRole := range userRoles {
		// append user to users map
		user, ok := usersMap[userRole.User.ID]
		if !ok {
			user = &userRole.User
			usersMap[user.ID] = user
		}

		usersMap[userRole.Role.UserID].Roles = append(usersMap[userRole.Role.UserID].Roles, &userRole.Role)
	}

	for _, u := range usersMap {
		user = *u
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": Payload{
			ID:    user.ID,
			Roles: user.Roles,
		},
		"exp": time.Now().Add(time.Minute * 15).Unix(),
	})

	tokenString, err := token.SignedString(JwtSecret)
	if err != nil {
		fmt.Println(err)
	}

	response := JwtResponse{
		Token: tokenString,
	}

	json.NewEncoder(w).Encode(response)
}
