package model

// RegisterClusterPayload is the payload received for registering cluster/runner with galaxy server
type RegisterClusterPayload struct {
	ClusterID  string `json:"id"`
	RunnerType string `json:"runnerType"`
	CreatedAt  string `json:"createdAt"`
	Url        string `json:"url"`
}

// LoginPayload is the payload received for user verification
type LoginPayload struct {
	Username string `json:"username"`
	Key      string `json:"key"`
}

type CrudRequestPayload struct {
	Op      string                 `json:"op"`
	Doc     interface{}            `json:"doc,omitempty"`
	Find    interface{}            `json:"find,omitempty"`
	Update  interface{}            `json:"update,omitempty"`
	Options map[string]interface{} `json:"options,omitempty"`
}
