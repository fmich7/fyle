package encryption

import "crypto/rand"

// GenerateRandomSalt generates a random salt
func GenerateRandomSalt() ([]byte, error) {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	return salt, err
}
