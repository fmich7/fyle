package server

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fmich7/fyle/pkg/config"
	"github.com/fmich7/fyle/pkg/storage"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleFileUpload(t *testing.T) {
	// create multipart form writer
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	_ = writer.WriteField("user", "testUser")
	_ = writer.WriteField("location", ".")

	fileWriter, err := writer.CreateFormFile("file", "testfile.txt")
	require.NoError(t, err, "Failed to create form file")
	_, err = fileWriter.Write([]byte("SOME DATA!!!"))
	require.NoError(t, err, "Failed to write data to file")

	err = writer.Close()
	assert.NoError(t, err, "Failed to close multipart writer")

	req := httptest.NewRequest(http.MethodPost, "/upload", &requestBody)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	recorder := httptest.NewRecorder()

	// create mock server instance with a mock store
	afs := afero.NewMemMapFs()
	storage, err := storage.NewTestingStorage(afs)
	if err != nil {
		t.Fatal(err)
	}
	mockServer := NewServer(config.NewTestingConfig(), storage)

	// call handler
	ctx := context.WithValue(req.Context(), CtxUsernameKey{}, "test")
	req = req.WithContext(ctx)

	mockServer.HandleFileUpload(recorder, req)

	// check response
	assert.Equal(t, http.StatusCreated, recorder.Code, "Expected HTTP 201 Created")
}
