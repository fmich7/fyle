package server_test

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fmich7/fyle/pkg/config"
	"github.com/fmich7/fyle/pkg/server"
	"github.com/fmich7/fyle/pkg/storage"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestHandleFileUpload(t *testing.T) {
	// create multipart form writer
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	_ = writer.WriteField("user", "testUser")
	_ = writer.WriteField("location", ".")

	fileWriter, err := writer.CreateFormFile("file", "testfile.txt")
	assert.NoError(t, err, "Failed to create form file")
	fileWriter.Write([]byte("SOME DATA!!!"))

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
	mockServer := server.NewServer(config.NewTestingConfig(), storage)

	// call handler
	ctx := context.WithValue(req.Context(), "username", "test")
	req = req.WithContext(ctx)

	mockServer.HandleFileUpload(recorder, req)

	// check response
	assert.Equal(t, http.StatusCreated, recorder.Code, "Expected HTTP 201 Created")
}
