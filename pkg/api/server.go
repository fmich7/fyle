package api

import (
	"net/http"

	"github.com/fmich7/fyle/pkg/storage"
)

type Server struct {
	listenAddr string
	store      storage.Storage
}

func NewServer(listenAddr string, store storage.Storage) *Server {
	return &Server{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *Server) Start() error {
	http.HandleFunc("POST /upload", s.HandleFileUpload)

	return http.ListenAndServe(s.listenAddr, nil)
}
