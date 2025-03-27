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
