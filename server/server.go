package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/launchpad/server/config"
	"github.com/spaceuptech/launchpad/utils/auth"
)

// Server modules manager the various clusters of launchpad
type Server struct {
	// For internal use
	router       *mux.Router
	config       *Config
	auth         *auth.Module
	galacyConfig *config.Module
}

// New creates a new launchpad server instance
func New(port, jwtPublicKeyPath, jwtPrivatePath, username, key string) (*Server, error) {

	// auth instance
	a := auth.Init(auth.Server, username, key)
	auth, err := auth.New(a, jwtPublicKeyPath, jwtPrivatePath)
	if err != nil {
		fmt.Errorf("error creating an instance of auth module - %v", err)
	}

	c, err := config.New(auth)
	if err != nil {
		fmt.Errorf("error creating an instance of config - %v", err)
	}

	return &Server{
		router:       mux.NewRouter(),
		config:       &Config{Port: port},
		auth:         auth,
		galacyConfig: c,
	}, nil
}

// Start begins the launchpad server operations
func (s *Server) Start() error {
	// Initialise the routes
	s.InitRoutes()

	// Start the launchpad server
	logrus.Infof("Starting launchpad server on port %s", s.config.Port)
	return http.ListenAndServe(":"+s.config.Port, s.router)
}
