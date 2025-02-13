package storage

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/fmich7/fyle/pkg/types"
)

type DiskStorage struct {
	location string
}

// NewDiskStorage creates a new DiskStorage object
func NewDiskStorage(fileUploadsLocation string) *DiskStorage {
	// Create the uploads directory if it doesn't exist
	if _, err := os.Stat(fileUploadsLocation); os.IsNotExist(err) {
		os.Mkdir(fileUploadsLocation, os.ModePerm)
	}

	return &DiskStorage{
		location: fileUploadsLocation,
	}
}

// UploadFile creates a file in the disk storage
func (d *DiskStorage) UploadFile(file *types.File) error {
	// create file in disk storage
	dst, err := d.createFile(file)
	if err != nil {
		return errors.New("Error creating file")
	}
	defer dst.Close()
	// write the file to disk
	if _, err := dst.ReadFrom(file.Data); err != nil {
		return errors.New("Error writing file to disk")
	}

	return nil
}

func (d *DiskStorage) DownloadFile(path string) error {
	return nil
}

// createFile creates a file in the disk storage
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

func (d *DiskStorage) GetFileUploadsLocation() string {
	return d.location
}
