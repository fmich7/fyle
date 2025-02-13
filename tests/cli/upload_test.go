package cli_test

import (
	"fmt"
	"io"

	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/fmich7/fyle/pkg/cli"
	"github.com/stretchr/testify/assert"
)

func TestUploadFile(t *testing.T) {
	// Create a temp file for testing
	tempFile, err := os.CreateTemp("", "testfile")
	assert.NoError(t, err, "Failed to create temporary file")
	defer os.Remove(tempFile.Name())

	// Write data to temp file
	_, err = tempFile.Write([]byte("SOME DATA!!!"))
	assert.NoError(t, err, "Failed to write content to temporary file")

	// Close temp file
	err = tempFile.Close()
	assert.NoError(t, err, "Failed to close temporary file")

	fmt.Println("Temp file path:", tempFile.Name()) // Debugging

	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify method
		assert.Equal(t, http.MethodPost, r.Method, "Expected POST request")

		// Verify content type
		assert.Contains(t, r.Header.Get("Content-Type"), "multipart/form-data", "Expected multipart form data")

		// Read request body
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err, "Failed to read request body")
		assert.Contains(t, string(body), "SOME DATA!!!", "Expected request body to contain file data")

		// Send response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("response body"))
	}))
	defer server.Close()

	// Set mock upload URL
	cli.UploadURL = server.URL

	// Call UploadFile
	err = cli.UploadFile(tempFile.Name(), "testLocation")

	// Assert no error
	assert.NoError(t, err, "Expected no error from UploadFile")
}
