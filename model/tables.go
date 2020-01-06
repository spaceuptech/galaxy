package model

// TableProjects is the structure of projects table
type TableProjects struct {
	ProjectID    string `json:"project_id,omitempty"`
	AccountID    string `json:"account_id,omitempty"`
	DefaultEnv   string `json:"default_environment,omitempty"`
	Environments string `json:"environments,omitempty"`
}

// TableServices is the structure of services table
type TableServices struct {
	ID          string `json:"id,omitempty"`
	Environment string `json:"environment,omitempty"`
	ProjectID   string `json:"project_id,omitempty"`
	ServiceID   string `json:"service_id,omitempty"`
	Clusters    string `json:"clusters,omitempty"` // csv
	Config      string `json:"config,omitempty"`   // json string
	Version     string `json:"version,omitempty"`
}

// TableClusters is the structure of clusters table
type TableClusters struct {
	ClusterID  string `json:"cluster_id,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	RunnerType string `json:"runner_type,omitempty"`
	Status     string `json:"status,omitempty"`
	Url        string `json:"url,omitempty"`
}
