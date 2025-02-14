package utils

import (
	"os"
	"path/filepath"

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

	// Get the absolute path of the file uploads location
	rootStoragePath, err := filepath.Abs(fileUploadsLocation)
	if err != nil {
		panic(err)
	}

	return &MockStore{
		location: rootStoragePath,
	}
}
