package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nsantiago2719/i_a_m/authorizer"
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

	var login LoginDetails
	if err := json.Unmarshal(body, &login); err != nil {
		fmt.Println("Failed unmarshal of login details: %w", err)
	}

	userRoles := []UserRole{}
	usersMap := make(map[string]*User)
	user := User{}
	query := `
  SELECT users.id AS "users.id",
       users.email AS "users.email",
       users.username AS "users.username",
       users.password AS "users.password",
       roles.name AS "roles.name",
       roles.permissions AS "roles.permissions",
       ur.user_id AS "userId"
  FROM users AS users
  LEFT JOIN users_roles AS ur ON users.id = ur.user_id
  LEFT JOIN roles AS roles ON roles.id = ur.role_id
  WHERE users.username=$1
  `

	if err := db.Select(&userRoles, query, login.Username); err != nil {
		fmt.Println("User select query failed: ", err)
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

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	claims := Claims{
		Payload{
			ID:    user.ID,
			Roles: user.Roles,
		},
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			ID:        uuid.NewString(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		fmt.Println(err)
	}

	response := JwtResponse{
		Token: tokenString,
	}

	json.NewEncoder(w).Encode(response)
}

// Logout handler blacklist the token if it doesnt exist
func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	redisClient := RedisClient()
	ctx := context.Background()
	// Get the token from the Authorization header
	token := r.Header.Get("Authorization")[7:]

	// Extract claim and checks validity
	claim, err := authorizer.ExtractClaim(token)
	if err != nil {
		response := GenericResponse{
			Message: "Token is invalid",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	// get the token from redis and returns a resault
	val, _ := redisClient.Get(ctx, claim.RegisteredClaims.ID).Result()

	// check if there is a value, if true, returns error
	if len(val) > 0 {
		errors.New("Token is already blacklisted")
		response := GenericResponse{
			Message: "Token is invalid",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Add token to cache for blacklisting
	// show error if there is any
	if err := redisClient.Set(ctx, claim.RegisteredClaims.ID, token, 15*time.Minute).Err(); err != nil {
		fmt.Println("Error: ", err)
	}

	response := GenericResponse{
		Message: "User logged out",
	}

	json.NewEncoder(w).Encode(response)
}
