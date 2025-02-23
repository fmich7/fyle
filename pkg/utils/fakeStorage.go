package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/fmich7/fyle/pkg/types"
	"github.com/spf13/afero"
)

type FakeStorage struct {
	AppFs    afero.Fs
	location string
}

// StoreFile always returns success
func (f *FakeStorage) StoreFile(file *types.File) error {
	// Get directory of file
	dirPath := filepath.Dir(file.Location)

	// Create path if it doesnt' exist
	if _, err := f.AppFs.Stat(dirPath); os.IsNotExist(err) {
		if err := f.AppFs.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("creating directory: %v\n", err)
		}
	}

	// Read file data as byte slice
	data, err := io.ReadAll(file.Data)
	if err != nil {
		return fmt.Errorf("reading file data: %v\n", err)
	}

	// Create file
	err = afero.WriteFile(f.AppFs, file.Location, data, 0644)
	if err != nil {
		return fmt.Errorf("writing file to a fake fs: %v\n", err)
	}

	return nil
}

// RetrieveFile always return success
func (f *FakeStorage) RetrieveFile(path string) (io.ReadCloser, error) {
	return f.AppFs.Open(path)
}

// GetFileUploadsLocation returns the location of the file uploads
func (m *FakeStorage) GetFileUploadsLocation() string {
	return m.location
}

// NewFakeStorage creates a new MockStore object
func NewFakeStorage(fileUploadsLocation string) (*FakeStorage, error) {
	// Get the absolute path of the file uploads location
	rootStoragePath, err := filepath.Abs(fileUploadsLocation)
	if err != nil {
		return nil, err
	}

	return &FakeStorage{
		AppFs:    afero.NewMemMapFs(),
		location: rootStoragePath,
	}, nil
}
