package auth_test

import (
	"testing"

	"github.com/fmich7/fyle/pkg/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var SECRET string = "hashingsecret"

func TestCreateToken(t *testing.T) {
	username := "admin"
	token, err := auth.CreateToken(SECRET, username)

	require.NoError(t, err)
	assert.Greater(t, len(token), 0)
}

func TestValidateToken(t *testing.T) {
	username := "admin"
	token, err := auth.CreateToken(SECRET, username)

	require.NoError(t, err)

	err = auth.ValidateToken(SECRET, token)
	assert.NoError(t, err)

	// invalid token
	invalidToken := "asdasd"
	err = auth.ValidateToken(SECRET, invalidToken)
	assert.Error(t, err)
}
