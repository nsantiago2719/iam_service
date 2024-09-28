package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nsantiago2719/iam_service/authorizer"
	"golang.org/x/crypto/bcrypt"
)

// Auth function is used for authentication of the user,
// returns a jwt containing the roles
// and authorized actions
func (s *API) handleAuth(w http.ResponseWriter, r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return APIError{
			Path:   "/auth",
			Status: http.StatusBadRequest,
			Msg:    err.Error(),
		}
	}

	var login LoginDetails
	if err := json.Unmarshal(body, &login); err != nil {
		return APIError{
			Path:   "/auth",
			Status: http.StatusBadRequest,
			Msg:    err.Error(),
		}
	}

	usersMap := make(map[string]*User)
	user := User{}
	permissions := []string{}

	userRoles, err := s.database.getUserWithRolesByUsername(login.Username)
	if err != nil {
		return APIError{
			Path:   "/auth",
			Status: http.StatusBadRequest,
			Msg:    err.Error(),
		}
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
		return APIError{
			Path:   "/auth",
			Status: http.StatusUnauthorized,
			Msg:    err.Error(),
		}
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return APIError{
			Path:   "/auth",
			Status: http.StatusUnauthorized,
			Msg:    err.Error(),
		}
	}

	response := JwtResponse{
		Token: tokenString,
	}

	JSONWriter(w, http.StatusOK, response)
	return nil
}

// Logout handler blacklist the token if it doesnt exist
func (s *API) handleLogout(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx := context.Background()
	// Get the token from the Authorization header
	token := r.Header.Get("Authorization")

	if token == "" {
		return APIError{
			Path:   "/logout",
			Status: http.StatusUnauthorized,
			Msg:    "Authorization header invalid",
		}
	}

	// Extract claim and checks validity
	sub, err := authorizer.AuthorizedAccess("user.logout", token[7:])
	if err != nil {
		return APIError{
			Path:   "/logout",
			Status: http.StatusUnauthorized,
			Msg:    err.Error(),
		}
	}
	// get the token from redis and returns a resault
	val, _ := s.memoryCache.Get(ctx, *sub).Result()

	// check if there is a value, if true, returns error
	if len(val) > 0 {
		return APIError{
			Path:   "/logout",
			Status: http.StatusBadRequest,
			Msg:    "Bad request",
		}
	}

	// Add token to cache for blacklisting
	// show error if there is any
	if err := s.memoryCache.Set(ctx, *sub, token, 15*time.Minute).Err(); err != nil {
		return APIError{
			Path:   "/logout",
			Status: http.StatusUnauthorized,
			Msg:    err.Error(),
		}
	}

	response := GenericResponse{
		Message: "User logged out",
	}

	JSONWriter(w, http.StatusOK, response)
	return nil
}
