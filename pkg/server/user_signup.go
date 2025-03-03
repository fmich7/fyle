package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fmich7/fyle/pkg/auth"
	"github.com/fmich7/fyle/pkg/types"
)

// CreateUser creates a new user and stores it in the database
func (s *Server) CreateUser(username, password string) error {
	// check if user already exists
	if _, err := s.store.RetrieveUser(username); err != nil {
		return err
	}

	acc, err := auth.NewUser(username, password)
	if err != nil {
		return fmt.Errorf("creating user: %v", err)
	}

	return s.store.StoreUser(acc)
}

// HandleSignUp handles sign up request
func (s *Server) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// get info from request
	var usrRequest types.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&usrRequest); err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	// Create new user
	if err := s.CreateUser(usrRequest.Username, usrRequest.Password); err != nil {
		http.Error(w, "Error creating an account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}
