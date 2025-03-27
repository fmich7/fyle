package cli

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewListTreeCmd(t *testing.T) {
	fs := afero.NewMemMapFs()
	client := NewCliClient(fs)

	cmd := client.NewListTreeCmd()
	require.NotNil(t, cmd)
	assert.Equal(t, "ls [serverPath]", cmd.Use)
}

func TestListFiles_Success(t *testing.T) {
	afs := afero.NewMemMapFs()
	tree := "|── ./\n|    |── test/\n|    |    |── file\n|    |    |── folder/\n"

	// mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(tree))
	}))
	defer server.Close()

	cli := NewCliClient(afs)
	cli.ListFilesURL = server.URL

	var buf bytes.Buffer
	cli.rootCmd.SetOut(&buf)
	err := cli.ListFiles("")

	assert.NoError(t, err, "Expected no error from ListFiles")
	assert.Equal(t, fmt.Sprintf("Your storage tree:\n%s", string(tree)), buf.String())
}
