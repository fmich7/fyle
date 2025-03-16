package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"golang.org/x/crypto/argon2"
)

// GeneratePBEKey generates a 32-byte encryption key for PBE with salt and KDF
// This key is used to encrypt/decrypt user data uploaded to the server
// The key itself is stored in the keyring
func GeneratePBEKey(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
}

// GenerateRandomNBytes generates a random byte slice of length n
func GenerateRandomNBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	return bytes, err
}

// EncryptData encrypts data using AES-GCM with nonce
func EncryptData(data io.Reader, encryptionKey []byte) (io.Reader, error) {
	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()

		// Create AES cipher block
		block, err := aes.NewCipher(encryptionKey)
		if err != nil {
			pw.CloseWithError(err)
			return
		}

		// Create GCM block
		gcm, err := cipher.NewGCM(block)
		if err != nil {
			pw.CloseWithError(err)
			return
		}

		// Generate random nonce
		nonce, err := GenerateRandomNBytes(gcm.NonceSize())
		if err != nil {
			pw.CloseWithError(err)
			return
		}

		if _, err := pw.Write(nonce); err != nil {
			pw.CloseWithError(err)
			return
		}

		// Write in chunks
		buffer := make([]byte, 4096)
		for {
			n, err := data.Read(buffer)
			if err == io.EOF {
				break
			} else if err != nil {
				pw.CloseWithError(err)
				return
			}

			cipherText := gcm.Seal(nil, nonce, buffer[:n], nil)
			if _, err := pw.Write(cipherText); err != nil {
				pw.CloseWithError(err)
				return
			}
		}
	}()

	return pr, nil
}
