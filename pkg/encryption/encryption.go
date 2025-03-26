package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
)

// GeneratePBEKey generates a 32-byte encryption key for PBE with salt.
// This key is used to encrypt/decrypt user files uploaded to the server.
// The key itself is stored in the keyring.
func GeneratePBEKey(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
}

// GenerateRandomNBytes generates a random byte slice of length n.
func GenerateRandomNBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	return bytes, err
}

// EncryptData encrypts data using AES-GCM with nonce (4KB chunks).
func EncryptData(data io.Reader, encryptionKey []byte) io.Reader {
	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()

		block, err := aes.NewCipher(encryptionKey)
		if err != nil {
			pw.CloseWithError(err)
			return
		}

		gcm, err := cipher.NewGCM(block)
		if err != nil {
			pw.CloseWithError(err)
			return
		}

		// generate nonce and write it to stream
		nonce := make([]byte, gcm.NonceSize())
		if _, err := rand.Read(nonce); err != nil {
			pw.CloseWithError(err)
			return
		}

		if _, err := pw.Write(nonce); err != nil {
			pw.CloseWithError(err)
			return
		}

		// 4KB chubnks
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

	return pr
}

// DecryptData decrypts data using AES-GCM with nonce (4KB size).
func DecryptData(data io.Reader, decryptionKey []byte) io.Reader {
	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()

		block, err := aes.NewCipher(decryptionKey)
		if err != nil {
			pw.CloseWithError(err)
			return
		}

		gcm, err := cipher.NewGCM(block)
		if err != nil {
			pw.CloseWithError(err)
			return
		}

		nonce := make([]byte, gcm.NonceSize())
		if _, err := io.ReadFull(data, nonce); err != nil {
			pw.CloseWithError(fmt.Errorf("failed to read nonce: %v", err))
			return
		}

		// 4kb + gcm chunks
		chunkSize := 4096 + gcm.Overhead()
		buffer := make([]byte, chunkSize)
		for {
			n, err := io.ReadFull(data, buffer)
			if err == io.EOF {
				break
			} else if err != nil && err != io.ErrUnexpectedEOF {
				pw.CloseWithError(fmt.Errorf("failed to read encrypted chunk: %v", err))
				return
			}

			decryptedChunk, err := gcm.Open(nil, nonce, buffer[:n], nil)
			if err != nil {
				pw.CloseWithError(fmt.Errorf("decryption failed: %v", err))
				return
			}

			if _, err := pw.Write(decryptedChunk); err != nil {
				pw.CloseWithError(fmt.Errorf("failed to write decrypted chunk: %v", err))
				return
			}
		}
	}()

	return pr
}
