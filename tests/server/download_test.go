package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/fmich7/fyle/pkg/server"
	"github.com/fmich7/fyle/pkg/storage"
	"github.com/fmich7/fyle/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleFileDownload(t *testing.T) {
	assert := assert.New(t)

	// create storage and mock file
	tmpDir := t.TempDir()
	filename := "tempfile.txt"
	user := "user"
	absUserDir := filepath.Join(tmpDir, user)
	absFilePath := filepath.Join(absUserDir, filename)
	content := []byte("some content")

	// create user dir
	err := os.MkdirAll(absUserDir, 0777)
	assert.NoError(err, "Expected no error creating user folder in db")

	// mock file
	file, err := os.Create(absFilePath)
	assert.NoError(err, "Expected no error on file creation")
	assert.FileExists(absFilePath, "Expected %s to be created", filename)

	// write data to file
	_, err = file.Write(content)
	file.Close()
	assert.NoError(err, "Expected no error on writing to %s", filename)

	// request
	body := new(bytes.Buffer)
	err = json.NewEncoder(body).Encode(types.DownloadRequest{
		Path: filename,
		User: user,
	})
	assert.NoError(err, "Expected no error with marshalling data")

	// server
	storage, err := storage.NewDiskStorage(tmpDir)
	require.NoError(t, err, "Expected no error creating storage")
	server := server.NewServer(":0", storage)

	// TEST: send request and validate if file matches
	req := httptest.NewRequest("POST", "/download", body)
	req.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	server.HandleFileDownload(response, req)

	data := response.Body.Bytes()
	assert.Equal(http.StatusOK, response.Code, string(data))
	assert.Equal(content, data, "Expected file content to be the same")

	// TEST: send request to get file that does not exists
	err = json.NewEncoder(body).Encode(types.DownloadRequest{
		Path: "DFHJADFKLJADJFDKLFLJKJKLF.txt",
		User: user,
	})
	assert.NoError(err, "Expected no error with marshalling data")

	req = httptest.NewRequest("POST", "/download", body)
	response = httptest.NewRecorder()
	server.HandleFileDownload(response, req)
	assert.Equal(http.StatusBadRequest, response.Code, string(data))

	// TEST: invalid path
	err = json.NewEncoder(body).Encode(types.DownloadRequest{
		Path: "../../../file.txt",
		User: user,
	})
	assert.NoError(err, "Expected no error with marshalling data")

	req = httptest.NewRequest("POST", "/download", body)
	response = httptest.NewRecorder()
	server.HandleFileDownload(response, req)
	assert.Equal(http.StatusBadRequest, response.Code, string(data))
}
