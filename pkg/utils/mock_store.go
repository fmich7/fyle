package utils

import (
	"os"

	"github.com/fmich7/fyle/pkg/types"
)

type MockStore struct {
	location string
}

// UploadFile always returns success
func (m *MockStore) UploadFile(file *types.File) error {
	return nil
}

// DownloadFile always return success
func (m *MockStore) DownloadFile(path string) error {
	return nil
}

// GetFileUploadsLocation returns the location of the file uploads
func (m *MockStore) GetFileUploadsLocation() string {
	return m.location
}

// NewMockStorage creates a new MockStore object
func NewMockStorage(fileUploadsLocation string) *MockStore {
	// Create the uploads directory if it doesn't exist
	if _, err := os.Stat(fileUploadsLocation); os.IsNotExist(err) {
		os.Mkdir(fileUploadsLocation, os.ModePerm)
	}

	return &MockStore{
		location: fileUploadsLocation,
	}
}

// Remove all files and directories from mockup db location
func (m *MockStore) Cleanup() {
	err := os.RemoveAll(m.location)
	if err != nil {
		panic(err)
	}
}
