package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var SECRET string = "hashingsecret"

func TestCreateToken(t *testing.T) {
	username := "admin"
	token, err := CreateToken(SECRET, username)

	require.NoError(t, err)
	assert.Greater(t, len(token), 0)
}

func TestValidateToken(t *testing.T) {
	username := "admin"
	token, err := CreateToken(SECRET, username)

	require.NoError(t, err)

	err = ValidateToken(SECRET, token)
	assert.NoError(t, err)

	// invalid token
	invalidToken := "asdasd"
	err = ValidateToken(SECRET, invalidToken)
	assert.Error(t, err)
}

func TestParseToken(t *testing.T) {
	username := "test"
	validToken, err := CreateToken(SECRET, username)
	require.NoError(t, err, "CreateToken should not return an error")

	_, err = ParseToken(SECRET, validToken)
	assert.NoError(t, err, "Error parsing token")

	invalidToken := validToken + "asd"
	_, err = ParseToken(SECRET, invalidToken)
	assert.Error(t, err)
}
