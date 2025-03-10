package cli_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fmich7/fyle/pkg/cli"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewUploadCmd ensures that the upload command is created successfully
func TestNewUploadCmd_ValidArgs(t *testing.T) {
	fs := afero.NewMemMapFs()
	client := cli.NewCliClient(fs)

	cmd := client.NewUploadCmd()
	require.NotNil(t, cmd)
	assert.Equal(t, "upload", cmd.Use)
}

// Mock server for handling file uploads
func mockUploadServer(expectedStatus int) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}
		w.WriteHeader(expectedStatus)
		w.Write([]byte(`{"message": "File uploaded successfully"}`))
	})
	return httptest.NewServer(handler)
}

// TestUploadFile_Success tests successful file upload
func TestUploadFile_Success(t *testing.T) {
	fs := afero.NewMemMapFs()
	client := cli.NewCliClient(fs)

	filePath := "/testfile.txt"
	afero.WriteFile(fs, filePath, []byte("test content"), 0644)

	server := mockUploadServer(http.StatusCreated)
	defer server.Close()
	client.UploadURL = server.URL

	err := client.UploadFile(filePath, "/server/path/")
	require.NoError(t, err)
}

// TestUploadFile_FailedRequest - failed HTTP request
func TestUploadFile_FailedRequest(t *testing.T) {
	fs := afero.NewMemMapFs()
	client := cli.NewCliClient(fs)

	filePath := "/testfile.txt"
	afero.WriteFile(fs, filePath, []byte("test content"), 0644)

	client.UploadURL = "http://invalid-url"

	err := client.UploadFile(filePath, "/server/path/")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "creating request")
}

// TestUploadFile_FormCreationError error in multipart form creation
func TestUploadFile_FormCreationError(t *testing.T) {
	fs := afero.NewMemMapFs()
	client := cli.NewCliClient(fs)

	err := client.UploadFile("/nonexistent.txt", "/server/path/")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "opening file")
}
