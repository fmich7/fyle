package api

import (
	"net/http"

	"github.com/fmich7/fyle/pkg/types"
)

type Server struct {
	listenAddr string
	store      types.Storage
}

func NewServer(listenAddr string, store types.Storage) *Server {
	return &Server{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *Server) Start() error {
	http.HandleFunc("POST /upload", s.HandleFileUpload)

	return http.ListenAndServe(s.listenAddr, nil)
}
