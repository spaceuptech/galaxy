package config

import (
	"context"
	"fmt"

	"github.com/spaceuptech/launchpad/model"
	"github.com/spaceuptech/launchpad/server/config/file"
	"github.com/spaceuptech/launchpad/utils"
)

type Module struct {
	block Config
}

type Config interface {
	AddProject(ctx context.Context, req *model.CreateProject) error
	GetProject(ctx context.Context, projectID string) (*model.CreateProject, error)
	GetProjects(ctx context.Context) ([]*model.CreateProject, error)
	DeleteProject(ctx context.Context, projectID string) error

	AddEnvironment(ctx context.Context, projectID string, req *model.Environment) error
	DeleteEnvironment(ctx context.Context, projectID, environmentId string) error

	AddCluster(ctx context.Context, projectID, environmentID string, req *model.Cluster) error
	DeleteCluster(ctx context.Context, projectID, environmentID, clusterID string) error
}

// New create a new instance of the Module object
func New(mode string) (*Module, error) {
	v, err := initBlock(mode)
	if err != nil {
		return nil, err
	}
	return &Module{block: v}, nil
}

func initBlock(mode string) (Config, error) {
	switch mode {
	case utils.CommunityEdition:
		return file.Init(), nil
	default:
		return nil, fmt.Errorf("error config invalid mode")
	}
}
