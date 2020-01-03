package model

// project configuration
type CreateProject struct {
	ID                 string         `json:"id" yaml:"id"`
	DefaultEnvironment string         `json:"default_environment" yaml:"default_environment"`
	Environments       []*Environment `json:"environment" yaml:"environment"`
}

type Environment struct {
	ID       string     `json:"id" yaml:"id"`
	Name     string     `json:"name" yaml:"name"`
	Clusters []*Cluster `json:"clusters" yaml:"clusters"`
}

type Cluster struct {
	ID  string `json:"id" yaml:"id"`
	Url string `json:"url" yaml:"url"`
}

type RegisterClusterRequest struct {
	ClusterID  string `json:"cluster_id"`
	RunnerType string `json:"runner_type"`
	CreatedAt  string `json:"created_at"`
	Url        string `json:"url"`
}

type InsertRequest struct {
	Query     string      `json:"query"`
	Variables interface{} `json:"variables"`
}

type MutationQueryResponse struct {
	Data struct {
		Update struct {
			Status int `json:"status"`
		} `json:"update,omitempty"`
		Insert struct {
			Status int `json:"status"`
		} `json:"insert,omitempty"`
	} `json:"data"`
}

type CliLoginRequest struct {
	Username string `json:"username"`
	Key      string `json:"key"`
}
