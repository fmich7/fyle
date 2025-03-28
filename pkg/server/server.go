package server

import (
	"fmt"
	"net"
	"net/http"
	"sync/atomic"

	"github.com/fmich7/fyle/pkg/config"
	"github.com/fmich7/fyle/pkg/storage"
)

// Server is a struct that represents the server.
type Server struct {
	listenAddr   string
	listener     atomic.Value
	store        storage.Storage
	jwtSecretKey string
	isRunning    atomic.Bool
}

// NewServer creates a new instance of the Server struct.
func NewServer(cfg *config.Config, store storage.Storage) *Server {
	return &Server{
		listenAddr:   cfg.ServerPort,
		store:        store,
		jwtSecretKey: cfg.JWTsecretKey,
		isRunning:    atomic.Bool{},
	}
}

// Start starts the server.
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return fmt.Errorf("failed to start listener: %v", err)
	}

	s.listener.Store(listener)

	mux := http.NewServeMux()
	// File related routes
	mux.HandleFunc("POST /file", s.AuthMiddleware(s.HandleFileUpload))
	mux.HandleFunc("POST /getfile", s.AuthMiddleware(s.HandleFileDownload))
	mux.HandleFunc("POST /ls", s.AuthMiddleware(s.HandleListFiles))

	// User related routes
	mux.HandleFunc("POST /signup", s.HandleSignUp)
	mux.HandleFunc("GET /login", s.HandleLogin)

	s.isRunning.Store(true)

	fmt.Println("Server started on", listener.Addr().String())
	return http.Serve(listener, mux)
}

// Shutdown gracefuly shutdowns the server
func (s *Server) Shutdown() error {
	l := s.listener.Load()
	if l != nil {
		if err := l.(net.Listener).Close(); err != nil {
			return fmt.Errorf("failed to close listener: %w", err)
		}
	}
	return s.store.Shutdown()
}

// GetPort returns the actual port used by the server.
func (s *Server) GetPort() (int, error) {
	l := s.listener.Load()
	if l == nil {
		return 0, fmt.Errorf("server is not running")
	}
	addr := l.(net.Listener).Addr().(*net.TCPAddr)
	return addr.Port, nil
}

// IsRunning return true if server is running
func (s *Server) IsRunning() bool {
	return s.isRunning.Load()
}
