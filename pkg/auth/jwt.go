package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// AuthClaims stores user data in jwt token.
type AuthClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// CreateToken creates token for user that lasts 24h.
func CreateToken(secret, username string) (string, error) {
	claims := AuthClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

// ValidateToken validates jwt token.
func ValidateToken(secret, tokenString string) error {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

// ParseToken parses jwt token and returns user claims.
func ParseToken(secret, tokenString string) (*AuthClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&AuthClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*AuthClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
