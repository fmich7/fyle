package storage

import (
	"io"

	"github.com/fmich7/fyle/pkg/auth"
	"github.com/fmich7/fyle/pkg/types"
)

// Storage represents a storage interface for server
type Storage interface {
	FileStorage
	UserStorage
}

// FileStorage represents file related methods
type FileStorage interface {
	StoreFile(file *types.File) error
	RetrieveFile(path string) (io.ReadCloser, error)
	GetFileUploadsLocation() string
	GetUserFileTree(string, string) (string, error)
}

// UserStorage represents user related methods
type UserStorage interface {
	StoreUser(user *auth.User) error
	RetrieveUser(username string) (*auth.User, error)
}

// ServerStorage aggregates multiple storages that match Storage interface
type ServerStorage struct {
	FileStorage
	UserStorage
}

// NewServerStorage returns new server storage
func NewServerStorage(fileDB FileStorage, userDB UserStorage) *ServerStorage {
	return &ServerStorage{
		fileDB,
		userDB,
	}
}
