package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Server modules manager the various clusters of launchpad
type Server struct {
	// For internal use
	router *mux.Router
	config *Config
}

// New creates a new launchpad server instance
func New(config *Config) *Server {
	return &Server{
		router: mux.NewRouter(),
		config: config,
	}
}

// Start begins the launchpad server operations
func (s *Server) Start() error {
	// Initialise the routes
	s.routes()

	// Start the launchpad server
	logrus.Infof("Starting launchpad server on port %s", s.config.Port)
	return http.ListenAndServe(":"+s.config.Port, s.router)
}
