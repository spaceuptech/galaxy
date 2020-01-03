package file

import "github.com/spaceuptech/launchpad/model"

// Config holds the entire configuration
type Config struct {
	Projects []*model.CreateProject `json:"projects" yaml:"projects"` // The key here is the project id
}

type CrudRequestBody struct {
	Op      string                 `json:"op"`
	Doc     interface{}            `json:"doc,omitempty"`
	Find    interface{}            `json:"find,omitempty"`
	Options map[string]interface{} `json:"options,omitempty"`
}
