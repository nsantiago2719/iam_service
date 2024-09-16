package authorizer

import (
	"errors"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx/types"
)

// Role used for creating struct for role data from database
type role struct {
	Name        string         `json:"name"`
	Permissions types.JSONText `json:"permissions"`
}

// payload contains the payload for the jwt token
type payload struct {
	roles []*role `json:"roles"`
}

type Claims struct {
	Data payload `json:"data"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// ExtractClaim returns the jwt claim and bool for token validitiy
func ExtractClaim(token string) (*Claims, error) {
	tokenVal, err := jwt.ParseWithClaims(token,
		&Claims{},
		func(_ *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
	if err != nil {
		return nil, err
	}

	if claims, ok := tokenVal.Claims.(*Claims); ok && tokenVal.Valid {
		return claims, nil
	}

	// by default return invalid token
	return nil, errors.New("Invalid token")
}

func AuthorizedAccess(resource, action, token string) error {
	_, err := ExtractClaim(token)
	if err != nil {
		fmt.Println(err)
	}

	// roles := []role{}
	return nil
}
