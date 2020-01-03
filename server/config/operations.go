package config

import (
	"context"

	"github.com/spaceuptech/launchpad/model"
)

func (m *Module) AddProject(ctx context.Context, req *model.CreateProject) error {
	return m.block.AddProject(ctx, req)
}

func (m *Module) GetProject(ctx context.Context, projectID string) ([]*model.TableProjects, error) {
	return m.block.GetProject(ctx, projectID)
}

func (m *Module) GetProjects(ctx context.Context) ([]*model.TableProjects, error) {
	return m.block.GetProjects(ctx)
}

func (m *Module) DeleteProject(ctx context.Context, projectID string) error {
	return m.block.DeleteProject(ctx, projectID)
}

func (m *Module) AddEnvironment(ctx context.Context, projectID string, req *model.Environment) error {
	return m.block.AddEnvironment(ctx, projectID, req)
}

func (m *Module) DeleteEnvironment(ctx context.Context, projectID, environmentId string) error {
	return m.block.DeleteEnvironment(ctx, projectID, environmentId)
}

func (m *Module) AddProjectCluster(ctx context.Context, projectID, environmentID string, req *model.Cluster) error {
	return m.block.AddProjectCluster(ctx, projectID, environmentID, req)
}

func (m *Module) DeleteProjectCluster(ctx context.Context, projectID, environmentID, clusterID string) error {
	return m.block.DeleteProjectCluster(ctx, projectID, environmentID, clusterID)
}

func (m *Module) UpsertService(ctx context.Context, req *model.Service) error {
	return m.block.UpsertService(ctx, req)
}

func (m *Module) DeleteService(ctx context.Context, serviceID string) error {
	return m.block.DeleteService(ctx, serviceID)
}

func (m *Module) ApplyService(ctx context.Context, req *model.Service) error {
	return m.block.ApplyService(ctx, req)
}
