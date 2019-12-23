package server

import (
	"github.com/gorilla/mux"

	"github.com/spaceuptech/launchpad/server/handlers"
)

func (s *Server) InitRoutes() {
	s.routes(s.router)
}

func (s *Server) routes(router *mux.Router) {

	router.Methods("POST").Path("/v1/galaxy/create-project").HandlerFunc(handlers.HandleServiceCreation())
	router.Methods("POST").Path("/v1/galaxy/create-service").HandlerFunc(handlers.HandleProjectCreation())

	// route for registering a new cluster
	router.Methods("POST").Path("/v1/galaxy/register-cluster").HandlerFunc(handlers.HandleClusterRegistration(s.auth))

	// route for login
	router.Methods("POST").Path("/v1/galaxy/login").HandlerFunc(handlers.HandleCliLogin(s.auth))

}
