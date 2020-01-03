package model

type TableProjects struct {
	ProjectID    string `json:"project_id,omitempty"`
	AccountID    string `json:"account_id,omitempty"`
	DefaultEnv   string `json:"default_environment,omitempty"`
	Environments string `json:"environments,omitempty"`
}
