package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// CreateToken creates token for user that lasts 24h
func CreateToken(secret, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(secret)
}

// ValidateToken validates given token
func ValidateToken(secret, tokenString string) error {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
