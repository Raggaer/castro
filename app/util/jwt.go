package util

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"github.com/kataras/go-errors"
)

type CastroClaims struct {
	Token string
	CreatedAt int64
	jwt.StandardClaims
}

// CreateJWToken signs and returns a new json web token
// with the given data inside the claim map
func CreateJWToken(claims CastroClaims) (string, error) {
	// Time for token to expire
	expires := time.Now().Add(time.Hour)

	claims.ExpiresAt = expires.Unix()
	claims.Issuer = "Castro AAC"

	// Create token with the given claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token and get as string
	return token.SignedString([]byte(Config.Secret))
}

// ParseJWToken reads the given json web token and
// returns the data map
func ParseJWToken(tkn string) (*CastroClaims, bool, error) {
	token, err := jwt.ParseWithClaims(tkn, &CastroClaims{}, func(token *jwt.Token) (interface{}, error) {

		// Check token signature
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, errors.New("Unexpected signing method")
		}

		return []byte(Config.Secret), nil
	})

	// Return any errors
	if err != nil {

		// Convert error to JWT type
		customError := err.(*jwt.ValidationError)

		// If token is expired tell the user
		if customError.Errors == jwt.ValidationErrorExpired {

			return nil, true, nil
		}

		return nil, false, err
	}

	// Grab token claims
	if claims, ok := token.Claims.(*CastroClaims); ok && token.Valid {

		// Check valid issuer
		if claims.Issuer != "Castro AAC" {

			return nil, false, errors.New("Invalid token issuer")
		}

		return claims, false, nil
	}



	return nil, false, errors.New("Cannot get token claims")
}