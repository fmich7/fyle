package server

import (
	"net/http"

	"github.com/fmich7/fyle/pkg/config"
	"github.com/fmich7/fyle/pkg/storage"
)

// Server is a struct that represents the server
type Server struct {
	listenAddr   string
	store        storage.Storage
	jwtSecretKey string
}

// NewServer creates a new instance of the Server struct
func NewServer(cfg *config.Config, store storage.Storage) *Server {
	return &Server{
		listenAddr:   cfg.ServerPort,
		store:        store,
		jwtSecretKey: cfg.JWTsecretKey,
	}
}

// Start starts the server
func (s *Server) Start() error {
	// File related routes
	http.HandleFunc("POST /file", s.HandleFileUpload)
	http.HandleFunc("POST /getfile", s.HandleFileDownload)

	// User related routes
	http.HandleFunc("POST /signup", s.HandleSignUp)
	http.HandleFunc("POST /login", s.HandleLogin)

	return http.ListenAndServe(s.listenAddr, nil)
}
