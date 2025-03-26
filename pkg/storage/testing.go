package storage

import (
	"github.com/spf13/afero"
)

// NewTestingStorage return storage that can be easily used in testing.
func NewTestingStorage(fs afero.Fs) (*ServerStorage, error) {
	userStorage := NewFakeUserDB()
	fileStorage, err := NewDiskFileStorage("uploads", fs)
	if err != nil {
		return nil, err
	}

	return NewServerStorage(
		fileStorage,
		userStorage,
	), nil
}
