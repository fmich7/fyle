package auth_test

import (
	"testing"

	"github.com/fmich7/fyle/pkg/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	username := "testuser"
	password := "securepassword"

	user, err := auth.NewUser(username, password)

	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, username, user.Username)
	assert.NotEqual(t, password, user.Password, "Password should be hashed")
	assert.Greater(t, len(user.Password), 0, "Hashed password should not be empty")
}

func TestCheckPassword_Valid(t *testing.T) {
	password := "securepassword"
	user, err := auth.NewUser("testuser", password)
	require.NoError(t, err)

	err = auth.CheckPassword(user.Password, password)
	assert.NoError(t, err, "Correct password should pass validation")
}

func TestCheckPassword_Invalid(t *testing.T) {
	password := "securepassword"
	user, err := auth.NewUser("testuser", password)
	require.NoError(t, err)

	err = auth.CheckPassword(user.Password, "wrongpassword")
	assert.Error(t, err, "Incorrect password should fail validation")
}
