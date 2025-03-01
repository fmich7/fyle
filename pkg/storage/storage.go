package storage

import (
	"io"

	"github.com/fmich7/fyle/pkg/types"
)

// Storage represents a storage interface for server
type Storage interface {
	StoreFile(file *types.File) error
	RetrieveFile(path string) (io.ReadCloser, error)
	GetFileUploadsLocation() string
}
