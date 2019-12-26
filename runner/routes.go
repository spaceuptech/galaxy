package runner

func (runner *Runner) routes() {
	runner.router.Methods("POST").Path("/v1/launchpad/project").HandlerFunc(runner.handleCreateProject())
	runner.router.Methods("POST").Path("/v1/launchpad/service").HandlerFunc(runner.handleServiceRequest())
	runner.router.HandleFunc("/v1/launchpad/socket", runner.handleWebsocketRequest())
	runner.router.HandleFunc("/v1/launchpad/manageServices/database", runner.handleDatabaseService())
}
