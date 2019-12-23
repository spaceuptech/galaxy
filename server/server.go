package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/launchpad/utils/auth"
)

// Server modules manager the various clusters of launchpad
type Server struct {
	// For internal use
	router *mux.Router
	config *Config
	auth   *auth.Module
}

// New creates a new launchpad server instance
func New(config *Config, auth *auth.Module) *Server {
	return &Server{
		router: mux.NewRouter(),
		config: config,
		auth:   auth,
	}
}

// Start begins the launchpad server operations
func (s *Server) Start() error {
	// Initialise the routes
	s.InitRoutes()

	// Start the launchpad server
	logrus.Infof("Starting launchpad server on port %s", s.config.Port)
	return http.ListenAndServe(":"+s.config.Port, s.router)
}
