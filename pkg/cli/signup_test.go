package cli

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSignUPCmd(t *testing.T) {
	fs := afero.NewMemMapFs()
	client := NewCliClient(fs)

	cmd := client.NewSignUPCmd()
	require.NotNil(t, cmd)
	assert.Equal(t, "signup [username] [password]", cmd.Use)
}

// TODO: i have to implement user mem db
func TestSignUPUser(t *testing.T) {
	afs := afero.NewMemMapFs()

	// mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("no error :)"))
	}))
	defer server.Close()

	cli := NewCliClient(afs)
	cli.SignupURL = server.URL

	var buf bytes.Buffer
	cli.rootCmd.SetOut(&buf)
	err := cli.SignUPUser("asdfasdf", "sadfadsf")

	assert.NoError(t, err, "Expected no error from Login")
	assert.Equal(t, "Created new account successfully!", buf.String())

	err = cli.LoginUser("", "")
	assert.Error(t, err)
}
