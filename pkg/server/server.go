package server

import (
	"net/http"

	"github.com/fmich7/fyle/pkg/config"
	"github.com/fmich7/fyle/pkg/storage"
)

// Server is a struct that represents the server
type Server struct {
	listenAddr string
	store      storage.Storage
}

// NewServer creates a new instance of the Server struct
func NewServer(cfg *config.Config, store storage.Storage) *Server {
	return &Server{
		listenAddr: cfg.ServerPort,
		store:      store,
	}
}

// Start starts the server
func (s *Server) Start() error {
	http.HandleFunc("POST /file", s.HandleFileUpload)
	// TODO: Fix this
	http.HandleFunc("POST /getfile", s.HandleFileDownload)

	return http.ListenAndServe(s.listenAddr, nil)
}
