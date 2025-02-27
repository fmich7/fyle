package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/fmich7/fyle/pkg/config"
	"github.com/fmich7/fyle/pkg/server"
	"github.com/fmich7/fyle/pkg/storage"
	"github.com/fmich7/fyle/pkg/types"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleFileDownload(t *testing.T) {
	assert := assert.New(t)

	// create storage and mock file
	filename := "tempfile.txt"
	user := "user"
	content := []byte("some content")

	// storage
	afs := afero.NewMemMapFs()

	storage, err := storage.NewDiskStorage("uploads", afs)
	require.NoError(t, err, "Expected no error creating storage")

	fileServerPath := filepath.Join(storage.GetFileUploadsLocation(), user, filename)
	afero.WriteFile(afs, fileServerPath, content, 0777)

	// server
	server := server.NewServer(config.NewTestingConfig(), storage)
	t.Log(storage.GetFileUploadsLocation())
	// request
	body := new(bytes.Buffer)
	err = json.NewEncoder(body).Encode(types.DownloadRequest{
		Path: filename,
		User: user,
	})
	assert.NoError(err, "Expected no error with marshalling data")

	// TEST: send request and validate if file matches
	req := httptest.NewRequest("POST", "/getfile", body)
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

	req = httptest.NewRequest("POST", "/getfile", body)
	response = httptest.NewRecorder()
	server.HandleFileDownload(response, req)
	assert.Equal(http.StatusBadRequest, response.Code, string(data))

	// TEST: invalid path
	err = json.NewEncoder(body).Encode(types.DownloadRequest{
		Path: "../../../file.txt",
		User: user,
	})
	assert.NoError(err, "Expected no error with marshalling data")

	req = httptest.NewRequest("POST", "/getfile", body)
	response = httptest.NewRecorder()
	server.HandleFileDownload(response, req)
	assert.Equal(http.StatusBadRequest, response.Code, string(data))
}
