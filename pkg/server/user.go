package server

import (
	"fmt"

	"github.com/fmich7/fyle/pkg/auth"
)

func (s *Server) CreateUser(username, password string) error {
	acc, err := auth.NewUser(username, password)
	if err != nil {
		return fmt.Errorf("creating user: %v", err)
	}

	// return s.store.StoreUser(acc)
	_ = acc
	return nil
}
