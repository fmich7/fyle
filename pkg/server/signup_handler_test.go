package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fmich7/fyle/pkg/config"
	"github.com/fmich7/fyle/pkg/storage"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleSignUp(t *testing.T) {
	afs := afero.NewMemMapFs()
	db, err := storage.NewTestingStorage(afs)
	require.NoError(t, err, "initializing db")

	cfg := config.NewTestingConfig()
	mockServer := NewServer(cfg, db)

	// Invalid signup request
	requestBody, err := json.Marshal(AuthUserRequest{
		Username: "",
		Password: "password123",
	})
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(requestBody))
	require.NoError(t, err)

	rec := httptest.NewRecorder()

	mockServer.HandleSignUp(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "length of passed username/password is 0")
}
