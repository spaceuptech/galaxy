package config

import (
	"context"

	"github.com/spaceuptech/launchpad/model"
	"github.com/spaceuptech/launchpad/server/config/file"
	"github.com/spaceuptech/launchpad/utils/auth"
)

type Module struct {
	block Config
	auth  *auth.Module
}

type Config interface {
	AddProject(ctx context.Context, req *model.CreateProject) error
	GetProject(ctx context.Context, projectID string) ([]*model.TableProjects, error)
	GetProjects(ctx context.Context) ([]*model.TableProjects, error)
	DeleteProject(ctx context.Context, projectID string) error

	AddEnvironment(ctx context.Context, projectID string, req *model.Environment) error
	DeleteEnvironment(ctx context.Context, projectID, environmentId string) error

	AddCluster(ctx context.Context, projectID, environmentID string, req *model.Cluster) error
	DeleteCluster(ctx context.Context, projectID, environmentID, clusterID string) error
}

// New create a new instance of the Module object
func New(auth *auth.Module) (*Module, error) {
	v, err := file.Init()
	if err != nil {
		return nil, err
	}
	return &Module{block: v, auth: auth}, nil
}
