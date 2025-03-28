package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/fmich7/fyle/pkg/config"
	"github.com/fmich7/fyle/pkg/storage"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleFileDownload(t *testing.T) {
	assert := assert.New(t)

	// data and fake storage
	filename := "tempfile.txt"
	user := "user"
	content := []byte("some content")
	afs := afero.NewMemMapFs()
	storage, err := storage.NewTestingStorage(afs)
	require.NoError(t, err, "Expected no error creating storage")

	// create mock file
	fileServerPath := filepath.Join(storage.GetFileUploadsLocation(), user, filename)
	require.NoError(
		t,
		afero.WriteFile(afs, fileServerPath, content, 0777),
		"Expected no error writing file",
	)

	// create srv
	srv := NewServer(config.NewTestingConfig(), storage)

	// sendRequest with injected username
	sendRequest := func(t *testing.T, path string) *httptest.ResponseRecorder {
		t.Helper()
		body := new(bytes.Buffer)
		require.NoError(
			t,
			json.NewEncoder(body).Encode(DownloadRequest{Path: path}),
			"Expected no error marshalling request data",
		)

		req := httptest.NewRequest("POST", "/getfile", body)
		req.Header.Set("Content-Type", "application/json")

		// Inject username into request context
		ctx := context.WithValue(req.Context(), CtxUsernameKey{}, user)
		req = req.WithContext(ctx)

		rec := httptest.NewRecorder()
		srv.HandleFileDownload(rec, req)

		return rec
	}

	// Valid file download
	response := sendRequest(t, filename)
	assert.Equal(http.StatusOK, response.Code, response.Body.String())
	assert.Equal(content, response.Body.Bytes(), "Expected file content to be the same")

	// Request for non-existent file
	response = sendRequest(t, "DFHJADFKLJADJFDKLFLJKJKLF.txt")
	assert.Equal(http.StatusBadRequest, response.Code, "Expected bad request for missing file")

	// Invalid file path (path traversal attack)
	response = sendRequest(t, "../../../file.txt")
	assert.Equal(http.StatusBadRequest, response.Code, "Expected bad request for invalid path")
}
