package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/launchpad/server/handlers"
)

func (s *Server) InitRoutes() {
	s.routes(s.router)
}

func (s *Server) routes(router *mux.Router) {

	// route for registering a new cluster
	router.Methods(http.MethodPost).Path("/v1/galaxy/register-cluster").HandlerFunc(handlers.HandleClusterRegistration(s.auth))

	// route for login
	router.Methods(http.MethodPost).Path("/v1/galaxy/login").HandlerFunc(handlers.HandleLogin(s.auth))

	// routes for galaxy configuration
	// projects
	router.Methods(http.MethodPost).Path("/v1/galaxy/project/create").HandlerFunc(handlers.HandleAddProject(s.auth, s.galacyConfig))
	router.Methods(http.MethodGet).Path("/v1/galaxy/project/{projectID}").HandlerFunc(handlers.HandleGetProject(s.auth, s.galacyConfig))
	router.Methods(http.MethodGet).Path("/v1/galaxy/projects").HandlerFunc(handlers.HandleGetProjects(s.auth, s.galacyConfig))
	router.Methods(http.MethodDelete).Path("/v1/galaxy/project/{projectID}").HandlerFunc(handlers.HandleDeleteProject(s.auth, s.galacyConfig))
	// project clusters
	router.Methods(http.MethodPost).Path("/v1/galaxy/project/{projectID}/{environmentID}/{clusterID}").HandlerFunc(handlers.HandleAddProjectCluster(s.auth, s.galacyConfig))
	router.Methods(http.MethodDelete).Path("/v1/galaxy/project/{projectID}/{environmentID}/{clusterID}").HandlerFunc(handlers.HandleDeleteProjectCluster(s.auth, s.galacyConfig))
	// project environments
	router.Methods(http.MethodPost).Path("/v1/galaxy/project/{projectID}/{environmentID}").HandlerFunc(handlers.HandleAddEnvironment(s.auth, s.galacyConfig))
	router.Methods(http.MethodDelete).Path("/v1/galaxy/project/{projectID}/{environmentID}").HandlerFunc(handlers.HandleDeleteEnvironment(s.auth, s.galacyConfig))
	// route for service configuration TODO ACTUAL IMPLEMENTATION REMAINING WITH DATABASE make sure it apply is upsert op
	router.Methods(http.MethodPost).Path("/v1/galaxy/service/{serviceID}").HandlerFunc(handlers.HandleUpsertService(s.auth, s.galacyConfig))
	router.Methods(http.MethodDelete).Path("/v1/galaxy/service/{serviceID}").HandlerFunc(handlers.HandleDeleteService(s.auth, s.galacyConfig))
	// web hook for applying service todo implementation reamaingig
	router.Methods(http.MethodPost).Path("/v1/galaxy/service/apply").HandlerFunc(handlers.HandleApplyService(s.auth, s.galacyConfig))

	// TODO implementation should be event based ?
	// route for providing public key to runner / cluster
	router.Methods(http.MethodGet).Path("/v1/galaxy/runner/key").HandlerFunc(handlers.HandleProvidePublicKey(s.auth))
}
