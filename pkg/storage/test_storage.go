package storage

import "github.com/spf13/afero"

// NewTestingStorage return storage that can be easily used in testing
func NewTestingStorage(fs afero.Fs) (*ServerStorage, error) {
	fileStorage, err := NewDiskFileStorage("uploads", fs)
	if err != nil {
		return nil, err
	}

	userStorage := NewDBUserStorage()

	return NewServerStorage(
		fileStorage,
		userStorage,
	), nil
}
