package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fmich7/fyle/pkg/auth"
	"github.com/fmich7/fyle/pkg/types"
)

// LoginUser return jwt token if user exists and password is correct
func (s *Server) LoginUser(username, password string) (string, error) {
	usr, err := s.store.RetrieveUser(username)
	if err != nil {
		return "", err
	}

	// validate password
	if err = auth.CheckPassword(usr.Password, password); err != nil {
		return "", err
	}

	return auth.CreateToken(s.jwtSecretKey, username)
}

// HandleLogin handles login request and returns jwt token on success
func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// get info from request
	var usrRequest types.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&usrRequest); err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	// login user
	token, err := s.LoginUser(usrRequest.Username, usrRequest.Password)
	if err != nil {
		http.Error(w, "Error invalid credentials", http.StatusUnauthorized)
		return
	}

	// send back jwt token
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, token)
}
