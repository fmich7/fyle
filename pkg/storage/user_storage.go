package storage

import (
	"github.com/fmich7/fyle/pkg/auth"
)

type DBUserStorage struct{}

func NewDBUserStorage() *DBUserStorage {
	return &DBUserStorage{}
}

func (d *DBUserStorage) StoreUser(user *auth.User) error {
	return nil
}
