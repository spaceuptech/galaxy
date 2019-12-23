package model

type ProjectCreateRequest struct {
	ProjectID string   `json:"project_id"`
	Clusters  []string `json:"clusters"`
	Team      string   `json:"team"`
	AccountID string   `json:"account_id"`
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
	Pass string `json:"pass"`
}