package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/fmich7/fyle/pkg/types"
	"github.com/fmich7/fyle/pkg/utils"
	"github.com/spf13/afero"
)

// DiskFileStorage is a struct that implements the Storage interface
// It is used to store files on disk
type DiskFileStorage struct {
	fs       afero.Fs
	location string
}

// NewDiskFileStorage creates a new DiskStorage object
// fs (nil is standard filesystem), you can pass aerof.Fs object for testing
func NewDiskFileStorage(fileUploadsLocation string, fs afero.Fs) (*DiskFileStorage, error) {
	// Default to OS filesystem if none is provided
	if fs == nil {
		fs = afero.NewOsFs()
	}

	// Create the uploads directory if it doesn't exist
	if _, err := os.Stat(fileUploadsLocation); os.IsNotExist(err) {
		if err := fs.Mkdir(fileUploadsLocation, os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed to create uploads directory: %v", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to check uploads directory: %v", err)
	}

	// Get the absolute path of the file uploads location
	rootStoragePath, err := filepath.Abs(fileUploadsLocation)
	if err != nil {
		return nil, fmt.Errorf("getting abs file upload location: %v", err)
	}

	return &DiskFileStorage{
		location: rootStoragePath,
		fs:       fs,
	}, nil
}

// UploadFile creates a file in the disk storage
func (d *DiskFileStorage) StoreFile(file *types.File) error {
	defer file.Data.Close()

	// Create file in disk storage
	dst, err := d.createFile(file)
	if err != nil {
		return err
	}

	defer dst.Close()

	// Write the file content
	if _, err := io.Copy(dst, file.Data); err != nil {
		return fmt.Errorf("writing file to disk: %v", err)
	}

	return nil
}

// RetrieveFile returns io.ReaderCloser of stored file
// Allows to copy content from file without loading it to memory
func (d *DiskFileStorage) RetrieveFile(path string) (io.ReadCloser, error) {
	return d.fs.Open(path)
}

// createFile creates a file in the disk storage
// It creates the directory if it doesn't exist
func (d *DiskFileStorage) createFile(file *types.File) (afero.File, error) {
	dirPath := filepath.Dir(file.Location)

	// Check if the directory exists, and create it if it doesn't
	if exists, _ := afero.DirExists(d.fs, dirPath); !exists {
		if err := d.fs.MkdirAll(dirPath, 0755); err != nil {
			return nil, fmt.Errorf("creating directory: %v", err)
		}
	}

	// Create the file
	return d.fs.Create(file.Location)
}

// GetFileUploadsLocation returns the file uploads location
func (d *DiskFileStorage) GetFileUploadsLocation() string {
	return d.location
}

// GetUserFileTree returns string with user's files
func (d *DiskFileStorage) GetUserFileTree(username, path string) (string, error) {
	// Get absolute valid path
	fullPath := utils.JoinPathParts(d.GetFileUploadsLocation(), username, path)

	// Check if valid
	valid := utils.ValidatePath(d.GetFileUploadsLocation(), fullPath)
	if !valid {
		return "", fmt.Errorf("invalid path")
	}

	treeString, err := utils.GetDirTree(d.fs, fullPath)
	if err != nil {
		return "", fmt.Errorf("building tree string %v", err)
	}

	return treeString, nil
}
