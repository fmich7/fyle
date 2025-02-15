package storage

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/fmich7/fyle/pkg/types"
)

// DiskStorage is a struct that implements the Storage interface
// It is used to store files on disk
type DiskStorage struct {
	location string
}

// NewDiskStorage creates a new DiskStorage object
func NewDiskStorage(fileUploadsLocation string) *DiskStorage {
	// Create the uploads directory if it doesn't exist
	if _, err := os.Stat(fileUploadsLocation); os.IsNotExist(err) {
		os.Mkdir(fileUploadsLocation, os.ModePerm)
	}

	// Get the absolute path of the file uploads location
	rootStoragePath, err := filepath.Abs(fileUploadsLocation)
	if err != nil {
		panic(err)
	}

	log.Println("Root storage path:", rootStoragePath)
	return &DiskStorage{
		location: rootStoragePath,
	}
}

// UploadFile creates a file in the disk storage
func (d *DiskStorage) StoreFile(file *types.File) error {
	// Create file in disk storage
	dst, err := d.createFile(file)
	if err != nil {
		return errors.New("Error creating file")
	}
	defer dst.Close()

	// Write the file to disk
	if _, err := dst.ReadFrom(file.Data); err != nil {
		return errors.New("Error writing file to disk")
	}

	return nil
}

// RetrieveFile returns io.ReaderCloser of stored file
func (d *DiskStorage) RetrieveFile(path string) (io.ReadCloser, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return file, err
}

// createFile creates a file in the disk storage
// It creates the directory if it doesn't exist
func (d *DiskStorage) createFile(file *types.File) (*os.File, error) {
	dirPath := filepath.Dir(file.Location)

	// Check if the directory exists, and create it if it doesn't
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			log.Printf("Failed to create directory: %v", err)
			return nil, err
		}
	}

	// Create the file
	dst, err := os.Create(file.Location)
	if err != nil {
		return nil, err
	}

	return dst, nil
}

// GetFileUploadsLocation returns the file uploads location
func (d *DiskStorage) GetFileUploadsLocation() string {
	return d.location
}
