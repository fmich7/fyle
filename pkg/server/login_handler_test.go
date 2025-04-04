package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fmich7/fyle/pkg/auth"
	"github.com/fmich7/fyle/pkg/config"
	"github.com/fmich7/fyle/pkg/storage"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginUser(t *testing.T) {
	afs := afero.NewMemMapFs()
	db, err := storage.NewTestingStorage(afs)
	require.NoError(t, err, "initializing db")

	cfg := config.NewTestingConfig()
	mockServer := NewServer(cfg, db)

	username := "testuser"
	password := "password"

	newUser, err := auth.NewUser(username, password)
	require.NoError(t, err, "creating user")
	db.StoreUser(newUser)

	// good password
	token, err := mockServer.LoginUser(username, password)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// wrong password
	token, err = mockServer.LoginUser(username, "asdads")
	assert.Error(t, err)
	assert.Empty(t, token)

	// user not found
	token, err = mockServer.LoginUser("asdsad", password)
	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestHandleLogin(t *testing.T) {
	afs := afero.NewMemMapFs()
	db, err := storage.NewTestingStorage(afs)
	require.NoError(t, err, "initializing db")

	cfg := config.NewTestingConfig()
	mockServer := NewServer(cfg, db)

	// invalid credentials
	requestBody, err := json.Marshal(AuthUserRequest{
		Username: "testuser",
		Password: "wrongpassword",
	})
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
	require.NoError(t, err)

	rec := httptest.NewRecorder()

	mockServer.HandleLogin(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "error invalid credentials")

	// valid login request
	username := "testuser"
	password := "password"

	newUser, err := auth.NewUser(username, password)
	require.NoError(t, err, "creating user")
	err = db.StoreUser(newUser)
	require.NoError(t, err, "failed to store user in test DB")

	requestBodyValid, err := json.Marshal(AuthUserRequest{
		Username: username,
		Password: password,
	})
	require.NoError(t, err)

	reqValid, err := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBodyValid))
	require.NoError(t, err)
	rec = httptest.NewRecorder()

	mockServer.HandleLogin(rec, reqValid)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotEmpty(t, rec.Body.String())
}
