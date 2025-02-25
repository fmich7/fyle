package cli_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/fmich7/fyle/pkg/cli"
	"github.com/fmich7/fyle/pkg/types"
	"github.com/stretchr/testify/assert"
)

var downloadErrMsg = "download error"

func mockDownloadFile(serverPath, destination string) error {
	if serverPath == "invalid/path" {
		return errors.New(downloadErrMsg)
	}
	return nil
}

func TestDownloadCommand(t *testing.T) {
	pwd, err := os.Getwd()
	assert.NoError(t, err)

	var output bytes.Buffer
	cmd := cli.NewDownloadCmd(mockDownloadFile)
	cmd.SetOut(&output)
	cmd.SetErr(&output)

	// Success
	cmd.SetArgs([]string{"valid/path", ""})
	err = cmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, output.String(), fmt.Sprintf("File saved at %s\n", pwd))

	// Fail
	cmd.SetArgs([]string{"invalid/path", ""})
	err = cmd.Execute()
	assert.NoError(t, err) // no error in cmd execute
	assert.Contains(t, output.String(), downloadErrMsg)
}

func TestDownloadFile(t *testing.T) {
	assert := assert.New(t)

	filename := "testfile.txt"
	content := []byte("SOME DATA!!")

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// Verify method
		assert.Equal(http.MethodPost, r.Method, "Expected POST request")

		// Verify content type
		assert.Contains(r.Header.Get("Content-Type"), "application/json", "Expected application/json content type")

		// Read request body
		data := new(types.DownloadRequest)
		err := json.NewDecoder(r.Body).Decode(data)
		assert.NoError(err, "Failed on unmarshalling data")

		// Create a temp file for testing
		buf := bytes.NewBuffer(content)

		// Send back file
		w.Header().Set("Content-Disposition", fmt.Sprintf(
			"attachment; filename=%s", filename,
		))
		w.Header().Set("Content-Type", "application/octet-stream")

		if _, err := io.Copy(w, buf); err != nil {
			http.Error(w, "Error streaming file", http.StatusInternalServerError)
			return
		}
	}))
	defer server.Close()

	// Set mock upload URL
	cli.DownloadURL = server.URL

	// Create folder where downloaded file will be stored
	tmpDir := t.TempDir()
	defer os.RemoveAll(tmpDir)

	// Test DownloadFile
	err := cli.DownloadFile(filename, tmpDir)
	assert.NoError(err, "Expected no error from DownloadFile")

	// Check if file exists in that path and if content is the same
	path := filepath.Join(tmpDir, filename)
	_, err = os.Stat(path)
	assert.NoError(err, "Expected file from request to be saved in provided destionatnio")

	fileContent, err := os.ReadFile(path)
	assert.NoError(err, "Expected no error no file read")
	assert.Equal(content, fileContent, "File contents does not match")
}
