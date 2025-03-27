package storage

import (
	"testing"

	"github.com/fmich7/fyle/pkg/auth"
	"github.com/stretchr/testify/assert"
)

func TestFakeUserDB(t *testing.T) {
	db := NewFakeUserDB()

	user1 := &auth.User{
		ID:       0,
		Username: "test",
		Password: "pass",
	}

	err := db.StoreUser(user1)
	assert.NoError(t, err, "expected no error while storing user")

	retrievedUser, err := db.RetrieveUser("test")
	assert.NoError(t, err, "expected no error while retrieving user")
	assert.Equal(t, user1, retrievedUser, "retrieved user does not match stored user")

	// duplicatekk
	err = db.StoreUser(user1)
	assert.Error(t, err, "expected error for duplicate user")

	// none exsisten
	_, err = db.RetrieveUser("nonexistent")
	assert.Error(t, err, "expected error for non-existent user")

	// shutdonw
	err = db.Shutdown()
	assert.NoError(t, err)
}
