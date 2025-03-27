package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/fmich7/fyle/pkg/auth"
)

// HandleSignUp handles sign up request.
func (s *Server) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// get info from request
	var usrRequest AuthUserRequest
	if err := json.NewDecoder(r.Body).Decode(&usrRequest); err != nil {
		log.Println(err)
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	// create new user
	user, err := auth.NewUser(usrRequest.Username, usrRequest.Password)
	if err != nil {
		log.Println(err)
		http.Error(w, "length of passed username/password is 0", http.StatusBadRequest)
		return
	}

	// store user
	if err := s.store.StoreUser(user); err != nil {
		log.Println(err)
		http.Error(w, "failed to store user on a server", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}
