package utils

import (
	"github.com/fmich7/fyle/pkg/config"
	"github.com/fmich7/fyle/pkg/storage"
	"github.com/spf13/afero"
)

// NewTestingStorage return storage that can be easily used in testing
func NewTestingStorage(fs afero.Fs) (*storage.ServerStorage, error) {
	userStorage := storage.NewFakeUserDB()
	fileStorage, err := storage.NewDiskFileStorage("uploads", fs)
	if err != nil {
		return nil, err
	}

	return storage.NewServerStorage(
		fileStorage,
		userStorage,
	), nil
}

// NewTestingConfig return config that is used for testing
func NewTestingConfig() *config.Config {
	return &config.Config{
		ServerPort:      ":0",
		JWTsecretKey:    "supersecret",
		UploadsLocation: "/uploads",
	}
}
