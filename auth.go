package main

import (
	"context"
	"encoding/json"
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
func (s *API) Auth(w http.ResponseWriter, r *http.Request) {
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
	permissions := []string{}
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

	if err := s.database.Select(&userRoles, query, login.Username); err != nil {
		fmt.Println("User select query failed: ", err)
	}

	for _, userRole := range userRoles {
		// append user to users map
		user, ok := usersMap[userRole.User.ID]
		if !ok {
			user = &userRole.User
			usersMap[user.ID] = user
		}

		// append permissions from role
		permissions = append(permissions, userRole.Role.Permissions)
	}

	for _, u := range usersMap {
		user = *u
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// remove duplicate scopes
	scope := RemoveStringDuplicate(permissions)

	claims := Claims{
		Payload{
			Scope: scope,
		},
		jwt.RegisteredClaims{
			Subject:   user.ID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			ID:        uuid.NewString(),
		},
	}

	fmt.Println(claims)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		fmt.Println(err)
	}

	response := JwtResponse{
		Token: tokenString,
	}

	JSONWriter(w, http.StatusOK, response)
}

// Logout handler blacklist the token if it doesnt exist
func (s *API) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := context.Background()
	// Get the token from the Authorization header
	token := r.Header.Get("Authorization")[7:]

	// Extract claim and checks validity
	claim, err := authorizer.ExtractClaim(token)
	if err != nil {
		response := GenericResponse{
			Message: "Token is invalid",
		}
		JSONWriter(w, http.StatusUnauthorized, response)
		return
	}
	// get the token from redis and returns a resault
	val, _ := s.memoryCache.Get(ctx, claim.RegisteredClaims.ID).Result()

	// check if there is a value, if true, returns error
	if len(val) > 0 {
		response := GenericResponse{
			Message: "Token is invalid",
		}
		JSONWriter(w, http.StatusUnauthorized, response)
		return
	}

	// Add token to cache for blacklisting
	// show error if there is any
	if err := s.memoryCache.Set(ctx, claim.RegisteredClaims.ID, token, 15*time.Minute).Err(); err != nil {
		fmt.Println("Error: ", err)
	}

	response := GenericResponse{
		Message: "User logged out",
	}

	JSONWriter(w, http.StatusOK, response)
}
