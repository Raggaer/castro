package util

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"github.com/kataras/go-errors"
	"github.com/dchest/uniuri"
	"github.com/raggaer/castro/app/database"
	"github.com/raggaer/castro/app/models"
	"github.com/jinzhu/gorm"
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

// CreateUniqueToken returns a unique token of the table sessions
// starting with the given length
func CreateUniqueToken(len int) (string, error) {
	// Create first token
	token := uniuri.NewLen(len)

	// Hold query results
	s := models.Session{}

	// Number of failed attempts
	fail := 0

	// Loop until token is unique
	for {

		// If we failed 3 times increment token length
		if fail > 3 {

			fail = 0
			len++
		}

		// Check if token is unique
		if err := database.DB.Table("sessions").First(&s, "token = ?", token).Error; err != nil {

			// Check if record exists
			if err == gorm.ErrRecordNotFound {

				// Stop loop
				break
			} else {

				return "", err
			}

			// Create new token
			token = uniuri.NewLen(len)

			// Increment fail attempts
			fail++
		}
	}

	return token, nil
}