package authorizer

import (
	"errors"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// payload contains the payload for the jwt token
type payload struct {
	Scope string `json:"scope"`
}

type Claims struct {
	Data payload `json:"data"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// ExtractClaim returns the jwt claim and bool for token validitiy
func extractClaim(token string) (*Claims, error) {
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

// AuthorizedAccess validates and checks if token is authorized to call the endpoint
// returns subject and error
func AuthorizedAccess(token string) (string, error) {
	claim, err := extractClaim(token)
	if err != nil {
		fmt.Println(err)
	}

	// roles := []role{}
	return claim.RegisteredClaims.Subject, nil
}
