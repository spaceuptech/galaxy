package config

// Config holds the entire configuration
type Config struct {
	Projects []*Project `json:"projects" yaml:"projects"` // The key here is the project id
}

// Project holds the project level configuration
type Project struct {
	ID      string   `json:"id" yaml:"id"`
	Name    string   `json:"name" yaml:"name"`
	Modules *Modules `json:"modules" yaml:"modules"`
}

// Modules holds the config of all the modules of that environment
type Modules struct {
	Environment map[string]*Environment `json:"environment" yaml:"environment"`
}

type Environment struct {
	Clusters []*Cluster `json:"clusters" yaml:"clusters"`
}

type Cluster struct {
	ID string `json:"id" yaml:"id"`
}
