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
	ClusterDead = "dead"
	ClusterAlive = "alive"
)
