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
	AddProject(ctx context.Context, accountID string, req *model.CreateProject) error
	GetProject(ctx context.Context, projectID, accountID string) ([]model.TableProjects, error)
	GetProjects(ctx context.Context, accountID string) ([]model.TableProjects, error)
	DeleteProject(ctx context.Context, projectID, accountID string) error
	CreateProjectInClusters(ctx context.Context, req *model.CreateProject) error
	DeleteProjectFromClusters(ctx context.Context, req *model.CreateProject) error

	AddEnvironment(ctx context.Context, accountID, projectID string, req *model.Environment) error
	DeleteEnvironment(ctx context.Context, accountID, projectID, environmentId string) error
	SetDefaultEnvironment(ctx context.Context, accountID, projectID, defaultEnv string) error

	AddProjectCluster(ctx context.Context, accountID, projectID, environmentID string, req *model.Cluster) error
	DeleteProjectCluster(ctx context.Context, accountID, projectID, environmentID, clusterID string) error
	UpdateCluster(ctx context.Context, request *model.RegisterClusterPayload, status string) error

	UpsertService(ctx context.Context, req *model.Service) error
	DeleteService(ctx context.Context, projectID, environmentID, serviceID, version string) error
	ApplyServiceToClusters(ctx context.Context, req *model.Service) error
	DeleteServiceFromClusters(ctx context.Context, req *model.Service) error
}

// New create a new instance of the Module object
func New(auth *auth.Module) (*Module, error) {
	v, err := file.Init(auth.GetUserName())
	if err != nil {
		return nil, err
	}
	return &Module{block: v, auth: auth}, nil
}
