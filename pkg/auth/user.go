package auth

import (
	"encoding/base64"
	"errors"

	"github.com/fmich7/fyle/pkg/encryption"
	"golang.org/x/crypto/bcrypt"
)

// User type
type User struct {
	ID       int
	Username string
	Password string
	Salt     string
}

// NewUser creates new user with hashed password
func NewUser(username, password string) (*User, error) {
	if len(username) == 0 || len(password) == 0 {
		return nil, errors.New("empty credentials")
	}

	// hash password
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	// generate random salt
	salt, err := encryption.GenerateRandomNBytes(16)
	if err != nil {
		return nil, err
	}

	return &User{
		Username: username,
		Password: hashedPassword,
		Salt:     base64.StdEncoding.EncodeToString(salt),
	}, nil

}

// CheckPassword checks if user provided valid password
func CheckPassword(storedPassword, enteredPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(enteredPassword))
}

// hashPassword hashes given password
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
