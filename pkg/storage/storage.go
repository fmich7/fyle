package storage

import (
	"io"

	"github.com/fmich7/fyle/pkg/auth"
	"github.com/fmich7/fyle/pkg/file"
)

// Storage represents a storage interface for server.
type Storage interface {
	FileStorage
	UserStorage
	Shutdown() error
}

// FileStorage represents file related methods.
type FileStorage interface {
	StoreFile(file *file.File) error
	RetrieveFile(path string) (io.ReadCloser, error)
	GetFileUploadsLocation() string
	GetUserFileTree(string) (string, error)
	Shutdown() error
}

// UserStorage represents user related methods.
type UserStorage interface {
	StoreUser(user *auth.User) error
	RetrieveUser(username string) (*auth.User, error)
	Shutdown() error
}

// ServerStorage aggregates multiple storages that match Storage interface.
type ServerStorage struct {
	FileStorage
	UserStorage
}

func (s *ServerStorage) Shutdown() error {
	if err := s.FileStorage.Shutdown(); err != nil {
		return err
	}
	return s.FileStorage.Shutdown()
}

// NewServerStorage returns new server storage.
func NewServerStorage(fileDB FileStorage, userDB UserStorage) *ServerStorage {
	return &ServerStorage{
		fileDB,
		userDB,
	}
}
