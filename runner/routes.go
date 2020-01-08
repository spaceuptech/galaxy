package runner

func (runner *Runner) routes() {
	runner.router.Methods("POST").Path("/v1/galaxy/project").HandlerFunc(runner.handleCreateProject())
	runner.router.Methods("POST").Path("/v1/galaxy/service").HandlerFunc(runner.handleServiceRequest())
	runner.router.HandleFunc("/v1/galaxy/socket", runner.handleWebsocketRequest())
	runner.router.HandleFunc("/v1/galaxy/manageServices/database", runner.handleDatabaseService())
}
