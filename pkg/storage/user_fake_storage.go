package storage

import "github.com/fmich7/fyle/pkg/auth"

type FakeUserDB struct{}

func NewFakeUserDB() *FakeUserDB {
	return &FakeUserDB{}
}

func (f *FakeUserDB) StoreUser(user *auth.User) error {
	return nil
}
func (f *FakeUserDB) RetrieveUser(username string) (*auth.User, error) {
	return nil, nil
}
