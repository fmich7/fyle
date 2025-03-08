package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/fmich7/fyle/pkg/auth"
	"github.com/fmich7/fyle/pkg/types"
)

// CreateUser creates a new user and stores it in the database
func (s *Server) CreateUser(username, password string) error {
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
		log.Println(err)
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	// Create new user
	if err := s.CreateUser(usrRequest.Username, usrRequest.Password); err != nil {
		log.Println(err)
		http.Error(w, "error creating an account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}
