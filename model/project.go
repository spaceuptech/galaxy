package model

// CreateProject describes the configuration of a project
type CreateProject struct {
	ID                 string         `json:"id" yaml:"id"`
	DefaultEnvironment string         `json:"defaultEnvironment" yaml:"defaultEnvironment"`
	Environments       []*Environment `json:"environment" yaml:"environment"`
}

// Environment is used to store environment information of particular project
type Environment struct {
	ID       string     `json:"id" yaml:"id"`
	Name     string     `json:"name" yaml:"name"`
	Clusters []*Cluster `json:"clusters" yaml:"clusters"`
}

// Cluster is used to store cluster information
type Cluster struct {
	ID  string `json:"id" yaml:"id"`
	Url string `json:"url" yaml:"url"`
}

// CreateClusterProjectPayload is the payload sent to cluster/runner for creating project
type CreateClusterProjectPayload struct {
	ProjectID    string        `json:"projectId"`
	Environments []Environment `json:"environments"`
}
