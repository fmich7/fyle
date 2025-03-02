package storage

import (
	"errors"

	"github.com/fmich7/fyle/pkg/auth"
)

type DBUserStorage struct{}

func NewDBUserStorage() *DBUserStorage {
	return &DBUserStorage{}
}

func (d *DBUserStorage) StoreUser(user *auth.User) error {
	return nil
}

func (d *DBUserStorage) RetrieveUser(username string) (*auth.User, error) {
	if username == "admin" {
		return nil, errors.New("that user exists already")
	}

	// just for now
	return auth.NewUser(username, "admin")
}
