package types

import "io"

// Storage represents a storage interface for server
type Storage interface {
	StoreFile(file *File) error
	RetrieveFile(path string) (io.ReadCloser, error)
	GetFileUploadsLocation() string
}
