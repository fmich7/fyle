package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/fmich7/fyle/pkg/auth"
)

// LoginUser return jwt token if user exists and password is correct.
func (s *Server) LoginUser(username, password string) (*LoginResponse, error) {
	usr, err := s.store.RetrieveUser(username)
	if err != nil {
		return nil, err
	}
	fmt.Printf("user %+v\n", usr)

	// validate password
	if err = auth.CheckPassword(usr.Password, password); err != nil {
		return nil, err
	}

	// generate token
	token, err := auth.CreateToken(s.jwtSecretKey, username)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		Salt:  usr.Salt,
	}, nil
}

// HandleLogin handles login request and returns jwt token on success.
func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// get info from request
	var usrRequest AuthUserRequest
	if err := json.NewDecoder(r.Body).Decode(&usrRequest); err != nil {
		log.Println(err)
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	// login user
	loginCredentials, err := s.LoginUser(usrRequest.Username, usrRequest.Password)
	if err != nil {
		log.Println(err)
		http.Error(w, "error invalid credentials", http.StatusUnauthorized)
		return
	}

	// send back token and salt
	jsonData, err := json.Marshal(LoginResponse{
		Token: loginCredentials.Token,
		Salt:  loginCredentials.Salt,
	})
	if err != nil {
		log.Println(err)
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
