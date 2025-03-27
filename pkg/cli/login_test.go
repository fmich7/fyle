package cli

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fmich7/fyle/pkg/server"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLoginCmd(t *testing.T) {
	fs := afero.NewMemMapFs()
	client := NewCliClient(fs)

	cmd := client.NewLoginCmd()
	require.NotNil(t, cmd)
	assert.Equal(t, "login [username] [password]", cmd.Use)
}

// TODO: i have to implement user mem db
func TestLoginUser(t *testing.T) {
	afs := afero.NewMemMapFs()

	// mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		credentials := server.LoginResponse{
			Token: "testToken",
			Salt:  "asdf",
		}

		data, err := json.Marshal(credentials)
		require.NoError(t, err)
		w.Write(data)
	}))
	defer server.Close()

	cli := NewCliClient(afs)
	cli.LoginURL = server.URL

	var buf bytes.Buffer
	cli.rootCmd.SetOut(&buf)
	err := cli.LoginUser("asdfasdf", "sadfadsf")

	assert.NoError(t, err, "Expected no error from Login")
	assert.Equal(t, "Logged in successfully!", buf.String())

	err = cli.LoginUser("", "")
	assert.Error(t, err)
}
