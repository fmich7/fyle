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

func TestNewUploadCmd_ValidArgs(t *testing.T) {
	fs := afero.NewMemMapFs()
	client := cli.NewCliClient(fs)

	cmd := client.NewUploadCmd()
	require.NotNil(t, cmd)
	assert.Equal(t, "upload [localPath] [serverPath]", cmd.Use)
}

// Mock server for handling file uploads
func mockUploadServer(expectedStatus int) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}
		w.WriteHeader(expectedStatus)
		w.Write([]byte("File uploaded successfully"))
	})
	return httptest.NewServer(handler)
}

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

func TestUploadFile_FailedRequest(t *testing.T) {
	fs := afero.NewMemMapFs()
	client := cli.NewCliClient(fs)

	filePath := "/testfile.txt"
	afero.WriteFile(fs, filePath, []byte("test content"), 0644)

	client.UploadURL = "http://invalid-url"

	err := client.UploadFile(filePath, "/server/path/")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "impossible to send a request")
}

func TestUploadFile_FormCreationError(t *testing.T) {
	fs := afero.NewMemMapFs()
	client := cli.NewCliClient(fs)

	server := mockUploadServer(http.StatusCreated)
	defer server.Close()
	client.UploadURL = server.URL

	err := client.UploadFile("/nonexistent.txt", "/server/path/")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "opening file")
}
