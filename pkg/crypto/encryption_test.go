package crypto

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptDecryptData(t *testing.T) {
	input := []byte("This is a secret message")
	plainReader := bytes.NewReader(input)

	encryptionKey, err := GenerateRandomNBytes(32)
	require.NoError(t, err)

	encryptedReader := EncryptData(plainReader, encryptionKey)
	var encryptedData bytes.Buffer
	_, err = io.Copy(&encryptedData, encryptedReader)
	require.NoError(t, err, "Encryption failed")

	decryptedReader := DecryptData(&encryptedData, encryptionKey)
	var decryptedData bytes.Buffer
	_, err = io.Copy(&decryptedData, decryptedReader)
	require.NoError(t, err, "Decryption failed")

	// check if input stayed the same
	assert.Equal(t, input, decryptedData.Bytes(), "Decrypted data does not match original")
}
