package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Server modules manager the various clusters of galaxy
type Server struct {
	// For internal use
	router *mux.Router
	config *Config
}

// New creates a new galaxy server instance
func New(config *Config) *Server {
	return &Server{
		router: mux.NewRouter(),
		config: config,
	}
}

// Start begins the galaxy server operations
func (s *Server) Start() error {
	// Initialise the routes
	s.routes()

	// Start the galaxy server
	logrus.Infof("Starting galaxy server on port %s", s.config.Port)
	return http.ListenAndServe(":"+s.config.Port, s.router)
}
