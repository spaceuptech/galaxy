package utils

const (
	// Graphql
	GraphqlEndpoint string = "http://localhost:4122/v1/api/space_galaxy/graphql"
	ApplicationJson string = "application/json"

	GraphqlQuery = iota
	GraphqlMutation

	// HTTP
	Ping = iota
	SimpleRequest

	MaximumPingRetries

	// Cluster
	ClusterDead  = "dead"
	ClusterAlive = "alive"

	// Config file
	ConfigFilePath = "config.yaml"

	// DB operations
	OpAll    = "all"
	OpOne    = "one"
	OpUpsert = "upsert"

	// SC crud operation endpoint
	CrudEndpoint = "http://localhost:4122/v1/api/spacegalaxy/crud/postgres"
	// SC crud operations
	OpCreate = "create"
	OpRead   = "read"
	OpUpdate = "update"
	OpDelete = "delete"

	// SC table names
	TableProjects = "projects"
	TableServices = "services"
	TableClusters = "clusters"
	TableAccounts = "accounts"

	// Sc project table fields
	ProjectID         = "project_id"
	ProjectAccount    = "account_id"
	ProjectDefaultEnv = "default_environment"
	ProjectEnvs       = "environments"
)
