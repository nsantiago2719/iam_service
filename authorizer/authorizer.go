package authorizer

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx/types"
)

// Role used for creating struct for role data from database
type Role struct {
	Name        string         `json:"name"`
	Permissions types.JSONText `json:"permissions"`
}

type Permission struct {
	Access   string   `json:"access"`
	Action   string   `json:"action"`
	Resource []string `json:"resource"`
}

// payload contains the payload for the jwt token
type Payload struct {
	Roles []*Role `json:"roles"`
}

type Claims struct {
	Data Payload `json:"data"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// GetValidClaim returns the jwt claim and bool for token validitiy
func GetValidClaim(token string) (*Claims, error) {
	tokenVal, err := jwt.ParseWithClaims(token,
		&Claims{},
		func(_ *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
	if err != nil {
		return nil, err
	}

	if claims, ok := tokenVal.Claims.(*Claims); ok && tokenVal.Valid {
		if auth := isAuthorized(claims); auth == true {
			return claims, nil
		}
	}

	// by default return invalid token
	return nil, errors.New("Invalid token")
}

func isAuthorized(claims *Claims) bool {
	permissions := []Permission{}
	for _, r := range claims.Data.Roles {
		json.Unmarshal(r.Permissions, &permissions)
	}

	fmt.Println(&permissions)

	return true
}
