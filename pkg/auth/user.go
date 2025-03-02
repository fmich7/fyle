package auth

import "golang.org/x/crypto/bcrypt"

// User type
type User struct {
	ID       int
	Username string
	Password string
}

// NewUser creates new user with hashed password
func NewUser(username, password string) (*User, error) {
	// hash password
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	// TODO: auto id
	return &User{
		ID:       0,
		Username: username,
		Password: hashedPassword,
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
