package storage

import (
	"io"

	"github.com/fmich7/fyle/pkg/auth"
	"github.com/fmich7/fyle/pkg/types"
)

// Storage represents a storage interface for server
type Storage interface {
	// File related methods
	StoreFile(file *types.File) error
	RetrieveFile(path string) (io.ReadCloser, error)
	GetFileUploadsLocation() string

	// User related methods
	StoreUser(user *auth.User) error
	RetrieveUser(username string) (*auth.User, error)
}

// ServerStorage aggregates multiple storages that match Storage interface
type ServerStorage struct {
	*DiskFileStorage
	*DBUserStorage
}

// NewServerStorage returns new server storage
func NewServerStorage(fileDB *DiskFileStorage, userDB *DBUserStorage) *ServerStorage {
	return &ServerStorage{
		DiskFileStorage: fileDB,
		DBUserStorage:   userDB,
	}
}
