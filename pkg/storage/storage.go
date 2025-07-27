package storage

import (
	"io"
	"sync"

	"github.com/fmich7/fyle/pkg/auth"
	"github.com/fmich7/fyle/pkg/config"
	"github.com/fmich7/fyle/pkg/file"
	"github.com/spf13/afero"
)

type Storage interface {
	RunMigrations(migrationsDir string) error
	Shutdown() error
}

// FileStorage represents file related methods.
type FileStorage interface {
	Storage
	StoreFile(file *file.File) error
	RetrieveFile(path string) (io.ReadCloser, error)
	GetFileUploadsLocation() string
	GetUserFileTree(string) (string, error)
}

// UserStorage represents user related methods.
type UserStorage interface {
	Storage
	StoreUser(user *auth.User) error
	RetrieveUser(username string) (*auth.User, error)
}

type MetadataStorage interface {
	Storage
	GetFileMetadata() error
	CreateFileMetadata() error
}

// ServerStorage aggregates multiple storages that match Storage interface.
type ServerStorage struct {
	FileStorage
	UserStorage
	MetadataStorage
}

// NewServerStorage returns new server storage.
func NewServerStorage(fileDB FileStorage, userDB UserStorage, metadataDB MetadataStorage) *ServerStorage {
	return &ServerStorage{
		fileDB,
		userDB,
		metadataDB,
	}
}

func (s *ServerStorage) Shutdown() []error {
	storages := []Storage{
		s.FileStorage,
		s.UserStorage,
		s.MetadataStorage,
	}
	errs := []error{}
	errCh := make(chan error, len(storages))
	var wg sync.WaitGroup

	wg.Add(len(storages))
	for _, storage := range storages {
		go func(s Storage) {
			defer wg.Done()
			errCh <- s.Shutdown()
		}(storage)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func (s *ServerStorage) SetUpStorages(cfg *config.Config) error {
	userStorage, err := NewPQUserStorage(cfg.UserPQCredentials)
	if err != nil {
		return err
	}
	s.UserStorage = userStorage

	fileStorage, err := NewDiskFileStorage(cfg.UploadsLocation, afero.NewOsFs())
	if err != nil {
		return err
	}
	s.FileStorage = fileStorage

	metadataStorage, err := NewPQMetadataStorage(cfg.MetadataPQCredentials)
	if err != nil {
		return err
	}
	s.MetadataStorage = metadataStorage

	return nil
}

func (s *ServerStorage) RunMigrations(migrationsDir string) []error {
	storages := []Storage{
		s.FileStorage,
		s.UserStorage,
		s.MetadataStorage,
	}
	errs := []error{}
	errCh := make(chan error, len(storages))
	var wg sync.WaitGroup

	wg.Add(len(storages))
	for _, storage := range storages {
		go func(s Storage) {
			defer wg.Done()
			errCh <- s.RunMigrations(migrationsDir)
		}(storage)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}
