package server

import (
	"net/http"

	"github.com/fmich7/fyle/pkg/types"
)

// Server is a struct that represents the server
type Server struct {
	listenAddr string
	store      types.Storage
}

// NewServer creates a new instance of the Server struct
func NewServer(listenAddr string, store types.Storage) *Server {
	return &Server{
		listenAddr: listenAddr,
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
