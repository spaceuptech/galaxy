package config

import (
	"context"

	"github.com/spaceuptech/launchpad/model"
)

func (m *Module) AddProject(ctx context.Context, accountID string, req *model.CreateProject) error {
	return m.block.AddProject(ctx, accountID, req)
}

func (m *Module) GetProject(ctx context.Context, accountID string, projectID string) ([]model.TableProjects, error) {
	return m.block.GetProject(ctx, accountID, projectID)
}

func (m *Module) GetProjects(ctx context.Context, accountID string) ([]model.TableProjects, error) {
	return m.block.GetProjects(ctx, accountID)
}

func (m *Module) DeleteProject(ctx context.Context, accountID string, projectID string) error {
	return m.block.DeleteProject(ctx, accountID, projectID)
}

func (m *Module) SetDefaultEnvironment(ctx context.Context, accountID string, projectID, defaultEnv string) error {
	return m.block.SetDefaultEnvironment(ctx, accountID, projectID, defaultEnv)
}

func (m *Module) AddEnvironment(ctx context.Context, accountID string, projectID string, req *model.Environment) error {
	return m.block.AddEnvironment(ctx, accountID, projectID, req)
}

func (m *Module) DeleteEnvironment(ctx context.Context, accountID string, projectID, environmentId string) error {
	return m.block.DeleteEnvironment(ctx, accountID, projectID, environmentId)
}

func (m *Module) AddProjectCluster(ctx context.Context, accountID string, projectID, environmentID string, req *model.Cluster) error {
	return m.block.AddProjectCluster(ctx, accountID, projectID, environmentID, req)
}

func (m *Module) DeleteProjectCluster(ctx context.Context, accountID string, projectID, environmentID, clusterID string) error {
	return m.block.DeleteProjectCluster(ctx, accountID, projectID, environmentID, clusterID)
}

func (m *Module) UpsertService(ctx context.Context, req *model.Service) error {
	return m.block.UpsertService(ctx, req)
}

func (m *Module) DeleteService(ctx context.Context, projectID, environmentID, serviceID string) error {
	return m.block.DeleteService(ctx, projectID, environmentID, serviceID)
}

func (m *Module) ApplyServiceToCluster(ctx context.Context, req *model.Service) error {
	return m.block.ApplyServiceToClusters(ctx, req)
}
