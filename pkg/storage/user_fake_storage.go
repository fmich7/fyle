package storage

import (
	"errors"

	"github.com/fmich7/fyle/pkg/auth"
)

// FakeUserDB simulates an memory user db.
type FakeUserDB struct {
	users map[string]auth.User
}

// NewFakeUserDB initialize db.
func NewFakeUserDB() *FakeUserDB {
	return &FakeUserDB{
		users: make(map[string]auth.User),
	}
}

// StoreUser adds an user.
func (f *FakeUserDB) StoreUser(user *auth.User) error {
	if _, ok := f.users[user.Username]; ok {
		return errors.New("user already exists")
	}

	f.users[user.Username] = auth.User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
	}
	return nil
}

// RetrieveUser reutrns an user from memory.
func (f *FakeUserDB) RetrieveUser(username string) (*auth.User, error) {
	user, exists := f.users[username]
	if !exists {
		return nil, errors.New("user not found")
	}

	return &auth.User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
	}, nil
}
