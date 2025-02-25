package utils

import (
	"io"
	"path/filepath"

	"github.com/fmich7/fyle/pkg/types"
)

type MockStore struct {
	location string
}

// StoreFile always returns success
func (m *MockStore) StoreFile(file *types.File) error {
	return nil
}

// RetrieveFile always return success
func (m *MockStore) RetrieveFile(path string) (io.ReadCloser, error) {
	return nil, nil
}

// GetFileUploadsLocation returns the location of the file uploads
func (m *MockStore) GetFileUploadsLocation() string {
	return m.location
}

// NewMockStorage creates a new MockStore object
func NewMockStorage(fileUploadsLocation string) *MockStore {
	// Get the absolute path of the file uploads location
	rootStoragePath, err := filepath.Abs(fileUploadsLocation)
	if err != nil {
		panic(err)
	}

	return &MockStore{
		location: rootStoragePath,
	}
}
