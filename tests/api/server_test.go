package api_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fmich7/fyle/pkg/api"
	"github.com/fmich7/fyle/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestHandleFileUpload(t *testing.T) {
	// Create multipart form writer
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add user and location form fields
	_ = writer.WriteField("user", "testUser")
	_ = writer.WriteField("location", ".")

	// Create temporary file field
	fileWriter, err := writer.CreateFormFile("file", "testfile.txt")
	assert.NoError(t, err, "Failed to create form file")

	// Write file content
	fileWriter.Write([]byte("SOME DATA!!!"))

	// Close the writer to finalize the request
	err = writer.Close()
	assert.NoError(t, err, "Failed to close multipart writer")

	// Create HTTP request
	req := httptest.NewRequest(http.MethodPost, "/upload", &requestBody)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create response recorder
	recorder := httptest.NewRecorder()

	// Create mock server instance with a mock store
	storage := utils.NewMockStorage("uploads")
	mockServer := api.NewServer(":0", storage)

	// Call handler
	mockServer.HandleFileUpload(recorder, req)

	// Check response
	assert.Equal(t, http.StatusCreated, recorder.Code, "Expected HTTP 201 Created")
}
