package cli_test

import (
	"io"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fmich7/fyle/pkg/cli"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestUploadFile(t *testing.T) {
	assert := assert.New(t)
	afs := afero.NewMemMapFs()

	// Create a temp file for testing
	filename := "testfile"
	err := afero.WriteFile(afs, filename, []byte("SOME DATA!!!"), 0644)
	assert.NoError(err, "Failed to create temporary file")

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify method
		assert.Equal(http.MethodPost, r.Method, "Expected POST request")

		// Verify content type
		assert.Contains(r.Header.Get("Content-Type"), "multipart/form-data", "Expected multipart form data")

		// Read request body
		body, err := io.ReadAll(r.Body)
		assert.NoError(err, "Failed to read request body")
		assert.Contains(string(body), "SOME DATA!!!", "Expected request body to contain file data")

		// Send response
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("response body"))
	}))
	defer server.Close()

	// Set mock upload URL
	cli := cli.NewCliClient(afs)
	cli.UploadURL = server.URL

	// Test UploadFile
	err = cli.UploadFile(filename, "testLocation")

	assert.NoError(err, "Expected no error from UploadFile")
}
