package cli_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fmich7/fyle/pkg/cli"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewDownloadCmd ensures that the download command is created successfully
func TestNewDownloadCmd_ValidArgs(t *testing.T) {
	fs := afero.NewMemMapFs()
	client := cli.NewCliClient(fs)

	cmd := client.NewDownloadCmd()
	require.NotNil(t, cmd)
	assert.Equal(t, "download", cmd.Use)
}

func TestDownloadFile_Success(t *testing.T) {
	assert := assert.New(t)
	afs := afero.NewMemMapFs()

	// Create a temp file for testing
	filename := "testfile"
	err := afero.WriteFile(afs, filename, []byte("SOME DATA!!!"), 0777)
	require.NoError(t, err, "Failed to create temporary file")

	// Create a directory to simulate the download folder
	err = afs.Mkdir("download", 0777)
	require.NoError(t, err, "Failed to create download directory")

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		w.Header().Set("Content-Type", "application/octet-stream")

		// send file
		file, err := afs.Open(filename)
		if err != nil {
			http.Error(w, "Error opening file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// copy file to respo
		_, err = io.Copy(w, file)
		if err != nil {
			http.Error(w, "Error copying file content", http.StatusInternalServerError)
		}

	}))
	defer server.Close()

	cli := cli.NewCliClient(afs)
	cli.DownloadURL = server.URL

	err = cli.DownloadFile(filename, "download")
	assert.NoError(err, "Expected no error from DownloadFile")

	exists, _ := afero.Exists(afs, "download/"+filename)
	assert.True(exists, "File should exist after download")
}

// TestDownloadFile_FailedRequest - failed HTTP request
func TestDownloadFile_FailedRequest(t *testing.T) {
	fs := afero.NewMemMapFs()
	client := cli.NewCliClient(fs)
	client.DownloadURL = "http://invalid-url"

	err := client.DownloadFile("server/path", "/local/path/")
	assert.Error(t, err)
	assert.Equal(t, "impossible to send a request", err.Error())
}
