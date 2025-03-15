package encryption

import (
	"crypto/rand"

	"golang.org/x/crypto/argon2"
)

// GeneratePBEKey generates a 32-byte encryption key for PBE with salt and KDF
// This key is used to encrypt/decrypt user data uploaded to the server
// The key itself is stored in the keyring
func GeneratePBEKey(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
}

// GenerateRandomSalt generates a random salt
func GenerateRandomSalt() ([]byte, error) {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	return salt, err
}
